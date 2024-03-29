import { Message, MessageId } from "./message";
import { MessageHandlers } from "./handler";

// 发送消息
export const sendMessage = (id: MessageId, msg: Message): Promise<void> => {
  return new Promise<void>((resolve, reject) => {
    // 序列化消息
    const buffer: ArrayBuffer = marshalMessage(id, msg, 0)
    // 发送消息
    wx.sendSocketMessage({
      data: buffer,
      success: () => {
        console.log('======= 消息发送成功: ', id);
        resolve();
      },
      fail: (error: WechatMiniprogram.GeneralCallbackResult) => {
        console.error('======= 消息发送失败：', id, error);
        reject(error);
      }
    });
  });
}

const responseHandlers: Map<number, (response: Message) => void> = new Map();
let requestSeq: number = 0; // 小程序单线程执行，不需要考虑线程安全

// 发送rpc请求
export const rpcMessage = (id: MessageId, msg: Message): Promise<Message> => {
  return new Promise<Message>((resolve, reject) => {
    const seq: number = ++requestSeq;
    // 注册响应处理函数
    responseHandlers.set(seq, (response: Message) => {
      resolve(response);
    });
    // 发送请求
    // 序列化消息
    const buffer: ArrayBuffer = marshalMessage(id, msg, seq)
    // 发送消息
    wx.sendSocketMessage({
      data: buffer,
      success: () => {
        console.log('======= rpc消息发送成功: ', id, ', seq=', seq);
      },
      fail: (error: WechatMiniprogram.GeneralCallbackResult) => {
        console.error('======= rpc消息发送失败: ', id, ', seq=', seq, error);
        responseHandlers.delete(seq);
        reject(error);
      }
    });
  });
}

// 处理接收到的服务器消息
export const onHandleMessage = (buffer: ArrayBuffer): void => {
  try {
    const dataView = new DataView(buffer);
    const id = dataView.getUint16(0, true);
    const seq = dataView.getUint32(2, true);
    console.log('======= receive server message id=', id, ', seq=', seq);
    const jsonBytes = new Uint8Array(buffer, 6);
    let jsonString = '';
    for (let i = 0; i < jsonBytes.length; i++) {
      jsonString += String.fromCharCode(jsonBytes[i]);
    }
    console.log('======= receive server message body=', jsonString);
    const message: Message = JSON.parse(jsonString);

    // 是否rpc消息
    if (seq > 0) {
      const handle = responseHandlers.get(seq);
      if (handle) {
        handle(message);
        responseHandlers.delete(seq);
        return;
      }
      console.warn('canot found rpc handle, message id: ', id, ', seq: ', seq);
    }
    // 根据消息id分发处理
    if (MessageHandlers[id]) {
      MessageHandlers[id](message);
    } else {
      console.error('unknown message id: ', id, ', seq: ', seq);
    }
  } catch (error) {
    console.error('on handle message error: ', error);
  }
}

// 序列化要发送的消息
const marshalMessage = (id: MessageId, msg: Message, seq: number): ArrayBuffer => {
  // 将 JSON 对象转换为字符串
  const jsonString = JSON.stringify(msg);
  // 字符串转为ArrayBuffer
  const jsonBuffer = new ArrayBuffer(jsonString.length);
  const jsonBufferView = new Uint8Array(jsonBuffer);
  for (let i = 0; i < jsonString.length; i++) {
    jsonBufferView[i] = jsonString.charCodeAt(i);
  }
  // 创建要发送的二进制数据
  const buffer = new ArrayBuffer(2 + 4 + jsonBuffer.byteLength);
  const dataView = new DataView(buffer);
  dataView.setUint16(0, id, true);
  dataView.setUint32(2, seq, true);
  // 将 JSON 数据复制到 ArrayBuffer 中
  const bufferView = new Uint8Array(buffer, 6); // 从 ArrayBuffer 的第 6 个字节开始填充
  bufferView.set(jsonBufferView);
  return buffer;
}
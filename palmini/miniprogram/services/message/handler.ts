import { Message, MessageId } from "./message";

// 定义消息处理函数的类型
type MessageHandler = (message: Message) => void;

interface IMessageHandler {
  [id: number]: MessageHandler
}

const handleLogin: MessageHandler = (message: Message) => {
  console.log('======= 处理登录消息返回: ', message);
};

const handleError: MessageHandler = (message: Message) => {
  console.error('======= 错误处理消息: ', message);
};

// 消息处理
export const MessageHandlers: IMessageHandler = {
  [MessageId.Login]: handleLogin,
}

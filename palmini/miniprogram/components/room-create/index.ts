import { Message, MessageId } from "../../services/message/message";
import { rpcMessage } from "../../services/message/send";

// components/room-create/index.ts
interface IComponentData {
  roomName: string;
  roomDesc: string;
}

interface CreateRoomReq extends Message {
  name: string;
  desc: string;
}
interface CreateRoomResp extends Message {
  code: number;
  room_id: string;
}

Component({

  /**
   * 组件的属性列表
   */
  properties: {
    showModal: {
      type: Boolean,
      value: false
    }
  },

  /**
   * 组件的初始数据
   */
  data: <IComponentData> {
    roomName: '',
    roomDesc: ''
  },

  /**
   * 组件的方法列表
   */
  methods: {
    // 输入房间名
    inputRoomName(e) {
      this.setData({
        roomName: e.detail.value
      });
    },
    // 输入房间描述
    inputRoomDesc(e) {
      this.setData({
        roomDesc: e.detail.value
      });
    },
    // 点击取消按钮
    cancel() {
      this.setData({
        showModal: false
      });
      wx.navigateTo({ url: '/pages/home/home' });
    },
    // 点击确认按钮
    confirm: async function() {
      if (this.data.roomName.trim() == '') {
        wx.showToast({
          title: '房间名不能为空',
          image: '/images/icon_error.png',
          duration: 1000
        });
        return;
      }
      // 关闭弹窗
      this.setData({
        showModal: false
      });
      // 向服务器发送请求
      const createRoomReq: CreateRoomReq = {
        name: this.data.roomName,
        desc: this.data.roomDesc,
      };
      const createRoomResp = await rpcMessage(MessageId.CreateRoom, createRoomReq) as CreateRoomResp
      console.log('======= create room response: ', createRoomResp.room_id);
      wx.showToast({
        title: '创建房间成功',
        image: '/images/icon_success.png',
        duration: 1000
      });
      // 这里可以添加创建房间的逻辑
      wx.navigateTo({ url: '/pages/home/home' });
    }
  }
})
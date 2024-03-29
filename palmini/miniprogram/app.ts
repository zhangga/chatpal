// app.ts
import updateManager from './services/utils/update';
import { config } from './config/index';
import { Message, MessageId } from './services/message/message';
import { onHandleMessage, rpcMessage } from './services/message/send';

interface LoginReq extends Message {
  code: string;
}
interface LoginResp extends Message {
  code: number;
}

App<IAppOption>({
  globalData: {},

  onLaunch: async function() {
    try {
      // 展示本地存储能力
      // const logs = wx.getStorageSync('logs') || []
      // logs.unshift(Date.now())
      // wx.setStorageSync('logs', logs)

      wx.getUserProfile({
        desc: '用于完善个人资料', // 声明获取用户个人信息后的用途
        success: res => {
          console.log('======= 用户信息: ', res);
        },
        fail: err => {
          console.error('======= 拒绝授权', err);
        }
      });
      
      // 显示连接服务器中
      wx.showLoading({
        title: '连接服务器中...',
      });
      // 使用await等待服务器连接的异步操作
      await this.connectServer();
      // 等待微信登录
      const code: string = await this.wxLogin();
      // 发送 res.code 到后台换取 openId, sessionKey, unionId
      const loginReq: LoginReq = {
        code: code,
      };
      const loginResp = await rpcMessage(MessageId.Login, loginReq) as LoginResp;
      console.log('======= login response code=', loginResp.code);

      // 隐藏加载中提示
      wx.hideLoading();

      // 数据加载完成后，跳转到首页
      wx.switchTab({
        url: '/pages/home/home', // 首页的路径
      });
    } catch (err) {
      // 隐藏加载中提示
      wx.hideLoading();
      wx.showToast({
        title: '连接服务器失败',
        icon: 'none',
      });
      console.error('======= 启动失败: ', err);
    }
  },
  onShow() {
    // updateManager();
  },
  // 连接服务器
  connectServer: function(): Promise<void> {
    return new Promise((resolve, reject) => {
      // 创建ws连接
      wx.connectSocket({
        url: config.wsAddr,
        header: {'content-type': 'application/json'},
        success: res => {
          console.log('======= websocket connect create success: ', res);
        },
        fail: err => {
          console.error('======= websocket connect create failed: ', err);
          reject(err);
        }
      });
      // 连接建立
      wx.onSocketOpen(function () {
        console.log('======= websocket connect success');
        resolve();
      });
      wx.onSocketError(err => {
        console.error('======= websocket connect error: ', err);
        reject(err);
      });
      wx.onSocketClose(function () {
        console.log('======= websocket connect closed！');
      });
      // 监听接收到服务器的消息事件
      wx.onSocketMessage(res => {
        onHandleMessage(res.data as ArrayBuffer)
      });
    });
  },
  // 微信登录
  wxLogin: function(): Promise<string> {
    return new Promise((resolve, reject) => {
      wx.login({
        success: res => {
          resolve(res.code);
        },
        fail: err => {
          console.error("======= wx.login() error: ", err)
          reject(err);
        }
      });
    });
  },
})
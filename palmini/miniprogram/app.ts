// app.ts
import updateManager from './utils/update';

App<IAppOption>({
  globalData: {
  },
  onLaunch() {
    // 展示本地存储能力
    const logs = wx.getStorageSync('logs') || []
    logs.unshift(Date.now())
    wx.setStorageSync('logs', logs)

    // 登录
    wx.login({
      success: res => {
        console.log("res.code: "+res.code)
        // 发送 res.code 到后台换取 openId, sessionKey, unionId
      },
      fail: err => {
        console.log("login error: ", err)
      }
    })

    // 连接服务器
    wx.connectSocket({
      url: "ws://10.7.160.151:10008/ws",
      header:{
        'content-type': 'application/json'
      },
      success: res => {
        console.log("connect socket success: ", res)
      }
    })
  },
  onShow() {
    // updateManager();
  },
})
/// <reference path="./types/index.d.ts" />

interface ConnectResponse {
  data: any;
}

interface IAppOption {
  globalData: {
    userInfo?: WechatMiniprogram.UserInfo,
  }
  userInfoReadyCallback?: WechatMiniprogram.GetUserInfoSuccessCallback,
  connectServer: () => Promise<void>,
  wxLogin: () => Promise<string>,
}
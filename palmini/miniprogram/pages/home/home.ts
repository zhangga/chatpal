// pages/lobby/lobby.ts
import { Room } from '../../models/room';
import { fetchHome } from '../../services/home/home';
import Toast from 'tdesign-miniprogram/toast/index';

interface PageData {
  pageLoading: boolean;
  roomList: Room[];
}

Page<PageData>({

  /**
   * 页面的初始数据
   */
  data: {
    pageLoading: false,
    roomList: [],
  },

  privateData: {
    tabIndex: 0,
  },

  methods: {
    sendSocketMessage() {
      console.log("handle tap")
      wx.sendSocketMessage({
        data: "hello i am chatpal",
      })
    },
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad() {
    this.init();
  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady() {

  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow() {

  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide() {

  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload() {

  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh() {
    this.init();
  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom() {

  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage() {

  },

  init() {
    this.loadHomePage();
  },

  loadHomePage() {
    wx.stopPullDownRefresh();

    this.setData({
      pageLoading: true,
    });
    fetchHome().then(({ roomList }) => {
      console.log("======= fetch roomList: ", roomList)
      this.setData({
        roomList: roomList,
        pageLoading: false,
      });
    });
  },

  navToSearchPage() {
    wx.navigateTo({ url: '/pages/search/index' });
  },
})
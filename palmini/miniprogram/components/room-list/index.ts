import { Room } from "../../models/room";

// components/room-list/index.ts
Component({
  externalClasses: ['wr-class'],

  /**
   * 组件的属性列表
   */
  properties: {
    roomList: {
      type: Array as () => Room[],
      value: [],
    },
    id: {
      type: String,
      value: '',
      observer: (id: string) => {
        this.genIndependentID(id);
      },
    },
    thresholds: {
      type: Array as () => number[],
      value: [],
    },
  },

  /**
   * 组件的初始数据
   */
  data: {
    independentID: '',
  },

  lifetimes: {
    created: function() {
    },
    ready: function() {
      this.init();
    },
  },

  /**
   * 组件的方法列表
   */
  methods: {
    onClickGoods(e: any) {
      const { index } = e.currentTarget.dataset;
      this.triggerEvent('click', { ...e.detail, index });
    },

    onClickGoodsThumb(e: any) {
      const { index } = e.currentTarget.dataset;
      this.triggerEvent('thumb', { ...e.detail, index });
    },

    init: function() {
      this.genIndependentID(this.id || '');
    },

    genIndependentID(id: string) {
      if (id) {
        this.setData({ independentID: id });
      } else {
        this.setData({
          independentID: `room-list-${~~(Math.random() * 10 ** 8)}`,
        });
      }
    },
  },
})
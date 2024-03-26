import {config} from '../../config/index';
import {Room} from '../../models/room';

export function fetchHome(): Promise<{roomList: Room[]}> {
  if (config.useMock) {
    return mockFetchHome();
  }
  return new Promise((resolve) => {
    resolve('real api');
  });
}

function mockFetchHome() {
  const { delay } = require('../_utils/delay');
  const { genSwiperRoomList } = require('../../models/swiper');
  return delay(3000).then(() => {
    return {
      roomList: genSwiperRoomList()
    };
  });
}
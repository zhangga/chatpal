import {Room} from './room'

const rooms: Room[] = [
  {name: "room_1", image: 'https://cdn-we-retail.ym.tencent.com/tsr/home/v2/banner1.png'},
  {name: "room_2", image: 'https://cdn-we-retail.ym.tencent.com/tsr/home/v2/banner2.png'},
  {name: "room_3", image: 'https://cdn-we-retail.ym.tencent.com/tsr/home/v2/banner3.png'},
  {name: "room_4", image: 'https://cdn-we-retail.ym.tencent.com/tsr/home/v2/banner4.png'},
  {name: "room_5", image: 'https://cdn-we-retail.ym.tencent.com/tsr/home/v2/banner5.png'},
  {name: "room_6", image: 'https://cdn-we-retail.ym.tencent.com/tsr/home/v2/banner6.png'},
];

export function genSwiperRoomList(): Room[] {
  return rooms;
}
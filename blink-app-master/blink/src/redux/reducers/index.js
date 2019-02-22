import { combineReducers } from 'redux';
import RoomReducer from './RoomReducer'

export default combineReducers({
  //Insert Reducers Here
  roomName: RoomReducer
});

console.log(combineReducers)

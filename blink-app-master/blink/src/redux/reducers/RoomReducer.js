import ROOM_NAME from '../actions';

export default (state = null, actions) => {
  switch (actions.type) {
    case 'ROOM_NAME':
      return actions.payload;
    default:
      return state;
  }
};

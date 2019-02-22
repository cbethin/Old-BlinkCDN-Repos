// Import libraries for making a component
import React from 'react';
import { Text, View, Image } from 'react-native';

// Make a Component
const Header = (props) => {
  const { imgStyle, viewStyle } = styles;
  // Rendering Component
  return (
    <View style={viewStyle}>
      <Image
          style={imgStyle}
          source={require('../../img/blinkChat.png')}
      />
    </View>
  );
}

// Component Styling
const styles = {
  viewStyle: {
    backgroundColor: '#F8F8F8',
    justifyContent: 'center',
    alignItems: 'center',
    height: 60,
    paddingTop: '5%',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.2,
    elevation: 2,
    position: 'relative'
  },
  imgStyle: {
    width: '60%',
    height: '100%'
  }
};

// Make the component avaliable to other parts of the app
export default Header;

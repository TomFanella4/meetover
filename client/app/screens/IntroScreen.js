import React from 'react';
import { StyleSheet } from 'react-native';
import AppIntroSlider from 'react-native-app-intro-slider';

import Colors from '../constants/Colors';

class IntroScreen extends React.Component {
  slides = [
    {
      key: 'welcome',
      title: 'Welcome to MeetOver',
      titleStyle: styles.text,
      text: 'Let\'s start with a quick tutorial',
      textStyle: styles.text,
      image: require('../../assets/images/icon.png'),
      imageStyle: styles.image,
      backgroundColor: Colors.tintColor,
    },
    {
      key: 'matches',
      title: 'Matches',
      titleStyle: styles.text,
      text: 'Start searching to view and be found by nearby professionals',
      textStyle: styles.text,
      image: require('../../assets/images/matches.jpg'),
      imageStyle: styles.image,
      backgroundColor: Colors.tintColor,
    },
    {
      key: 'map',
      title: 'Map',
      titleStyle: styles.text,
      text: 'Browse professionals nearby',
      textStyle: styles.text,
      image: require('../../assets/images/map.jpg'),
      imageStyle: styles.image,
      backgroundColor: Colors.tintColor,
    },
    {
      key: 'meet',
      title: 'Meet',
      titleStyle: styles.text,
      text: 'Chat with professionals you\'d like to meet',
      textStyle: styles.text,
      image: require('../../assets/images/meet.jpg'),
      imageStyle: styles.image,
      backgroundColor: Colors.tintColor,
    },
    {
      key: 'chat',
      title: 'Chat',
      titleStyle: styles.text,
      text: 'Pick a location to meet and enjoy the conversation!',
      textStyle: styles.text,
      image: require('../../assets/images/chat.jpg'),
      imageStyle: styles.image,
      backgroundColor: Colors.tintColor,
    }
  ];

  _handleDone() {
    this.props.navigation.navigate('Login');
  }

  render() {
    return (
      <AppIntroSlider
        slides={this.slides}
        onDone={() => this._handleDone()}
      />
    );
  }
}

export default IntroScreen;

const styles = StyleSheet.create({
  image: {
    height: 400,
    resizeMode: 'contain'
  },
  text: {
    fontFamily: 'pt-sans'
  }
});

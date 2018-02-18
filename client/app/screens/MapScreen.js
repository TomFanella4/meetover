import React from 'react';
import { View, StyleSheet } from 'react-native';

import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText'

export default class MapScreen extends React.Component {
  static navigationOptions = {
    title: 'Map',
    headerStyle: {
      backgroundColor: Colors.tintColor,
    },
    headerTintColor: Colors.header,
  };

  render() {
    return (
      <View style={styles.container}>
        <PTSansText>Map Screen</PTSansText>
      </View>
    );
  }
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingTop: 15,
    backgroundColor: '#fff',
  }
});

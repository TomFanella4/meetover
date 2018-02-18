import React from 'react';
import { ScrollView, StyleSheet } from 'react-native';

import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText'

export default class ChatsScreen extends React.Component {
  static navigationOptions = {
    title: 'Chats',
    headerStyle: {
      backgroundColor: Colors.tintColor,
    },
    headerTintColor: Colors.header,
  };

  render() {
    return (
      <ScrollView style={styles.container}>
        <PTSansText>Chats Screen</PTSansText>
      </ScrollView>
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

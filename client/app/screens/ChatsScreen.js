import React from 'react';
import { StyleSheet } from 'react-native';
import { Container, Content, List, ListItem } from 'native-base';

import { PTSansText } from '../components/StyledText'
import IsSearchingBar from '../components/IsSearchingBar';

export default class ChatsScreen extends React.Component {
  static navigationOptions = {
    title: 'Chats',
  };

  _viewChatThread() {
    this.props.navigation.navigate('ChatScreen');
  }

  render() {
    return (
      <Container style={styles.container}>
        <IsSearchingBar />
        <Content>
          <List>
            <ListItem onPress={() => this._viewChatThread()}>
              <PTSansText>Krutarth Rao</PTSansText>
            </ListItem>
          </List>
        </Content>
      </Container>
    );
  }
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  }
});

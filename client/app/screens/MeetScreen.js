import React from 'react';
import { StyleSheet } from 'react-native';
import {
  Container,
  Content,
  View,
  List,
  ListItem,
  Left,
  Body,
  Spinner
} from 'native-base';
import { connect } from 'react-redux';

import { PTSansText } from '../components/StyledText';
import { ProfileImage } from '../components/ProfileImage';
import Colors from '../constants/Colors';
import { chatScreenStrings } from '../constants/Strings';
import IsSearchingBar from '../components/IsSearchingBar';

class MeetScreen extends React.Component {
  static navigationOptions = {
    title: 'Meet',
  };

  _viewThread(thread) {
    const { navigation } = this.props;

    switch(thread.status) {
      case 'pending':
        if (thread.origin === 'receiver') {
          navigation.navigate('ConfirmScreen', thread);
          return;
        }
        break;

      case 'accepted':
        navigation.navigate('ChatScreen', thread);
        return;

      case 'declined':
        if (thread.origin === 'receiver') {
          navigation.navigate('RequestScreen', thread);
          return;
        }
        break;
    }

    navigation.navigate('ProfileScreen', thread);
  }

  _buildThreadList(threads) {
    const acceptedItems = [], pendingItems = [], declinedItems = [], threadItems = [];
    threads.forEach(thread => {
      const threadItem = (
        <ListItem
          key={thread._id}
          onPress={() => this._viewThread(thread)}
          avatar
        >
          <Left>
            <ProfileImage pictureUrl={thread.profile.pictureUrl} />
          </Left>
          <Body>
            <PTSansText style={styles.name}>{thread.profile.formattedName}</PTSansText>
            <PTSansText note>{thread.profile.headline}</PTSansText>
          </Body>
        </ListItem>
      );

      switch(thread.status) {
        case 'pending':
          pendingItems.push(threadItem);
          break;

        case 'accepted':
          acceptedItems.push(threadItem);
          break;

        case 'declined':
          declinedItems.push(threadItem);
          break;
      }
    });

    if (pendingItems.length > 0) {
      threadItems.push(<ListItem key={'pending'} itemDivider><PTSansText>Pending</PTSansText></ListItem>);
      threadItems.push(...pendingItems);
    }

    if (acceptedItems.length > 0) {
      threadItems.push(<ListItem key={'accepted'} itemDivider><PTSansText>Accepted</PTSansText></ListItem>);
      threadItems.push(...acceptedItems);
    }

    if (declinedItems.length > 0) {
      threadItems.push(<ListItem key={'declined'} itemDivider><PTSansText>Declined</PTSansText></ListItem>);
      threadItems.push(...declinedItems);
    }

    return threadItems;
  }

  render() {
    const { threadList } = this.props;
    const threads = threadList ? Object.values(threadList) : null;
    let threadItems;

    if (threadList) {
      threadItems = this._buildThreadList(threads);
    }

    return (
      <Container style={styles.container}>
        <IsSearchingBar />
        {
          threadItems ?
            threadItems.length > 0 ?
              <Content>
                <List>
                  {threadItems}
                </List>
              </Content>
            :
              <View style={styles.noChatsText}>
                <PTSansText>{chatScreenStrings.noChatsText}</PTSansText>
              </View>
          :
            <Spinner color={Colors.tintColor} />
        }
      </Container>
    );
  }
};

const mapStateToProps = state => ({
  threadList: state.chat.threadList
});

export default connect(
  mapStateToProps
)(MeetScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  noChatsText: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center'
  },
  name: {
    fontSize: 20
  }
});

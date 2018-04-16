import React from 'react';
import { StyleSheet } from 'react-native';
import { connect } from 'react-redux';
import {
  Spinner,
  View,
  Button,
  Text
} from 'native-base';
import { GiftedChat, Bubble } from 'react-native-gifted-chat';

import MeetOverModal from '../components/MeetOverModal';
import { PTSansText } from '../components/StyledText';
import { ProfileImage } from '../components/ProfileImage';
import { chatMessagesToLoad } from '../constants/Common';
import Colors from '../constants/Colors';
import { sendFirebaseMessage, modifyFirebaseUserState } from '../firebase';
import {
  registerFetchNewMessageAsync,
  fetchEarlierMessagesAsync
} from '../actions/chatActions';

class ChatScreen extends React.Component {
  static navigationOptions = ({ navigation }) => ({
    title: navigation.state.params.profile.formattedName,
    headerRight: (
      <Button
        onPress={navigation.state.params._handleConclude}
        transparent
        style={styles.headerButton}
      >
        <Text style={styles.headerText}>Conclude</Text>
      </Button>
    )
  });

  state = {
    isModalVisible: false,
    modalViewType: !this.props.navigation.state.params.review ? 'conclude' : 'completed'
  }

  _handleSend(messages = []) {
    const _id = this.props.navigation.state.params._id;
    const access_token = this.props.userProfile.token.access_token;

    sendFirebaseMessage(_id, messages, access_token);
  }

  _handleLoadEarlier() {
    const {
      messages,
      navigation,
      fetchEarlierMessagesAsync,
    } = this.props;
    const _id = navigation.state.params._id;

    fetchEarlierMessagesAsync(_id, chatMessagesToLoad, messages[messages.length - 1]._id);
  }

  _handlePressAvatar() {
    const { navigation } = this.props;
    navigation.navigate('ProfileScreen', { profile: navigation.state.params.profile });
  }

  renderBubble(props) {
    return (
      <Bubble
        {...props}
        wrapperStyle={{
          right: {
            backgroundColor: Colors.tintColor
          }
        }}
      />
    )
  }

  async _handleSubmitReview(review) {
    const threadID = this.props.navigation.state.params._id;
    await modifyFirebaseUserState(`threadList/${threadID}/review`, review);
    this.setState({ modalViewType: 'completed' });
  }

  _renderConcludeModal() {
    return (
      <View>
        <PTSansText style={styles.modalHeaderText}>How was your MeetOver experience?</PTSansText>
        <Button style={styles.modalButton} onPress={() => this._handleSubmitReview('positive')}>
          <PTSansText style={styles.modalText}>It was a good match</PTSansText>
        </Button>
        <Button style={styles.modalButton} onPress={() => this._handleSubmitReview('negative')}>
          <PTSansText style={styles.modalText}>It wasn't a good match</PTSansText>
        </Button>
        <Button style={styles.modalButton} onPress={() => this._handleSubmitReview('none')}>
          <PTSansText style={styles.modalText}>We didn't meet</PTSansText>
        </Button>
      </View>
    );
  }

  _renderCompletedModal() {
    return (
      <View>
        <PTSansText style={styles.modalHeaderText}>Thanks for your feedback!</PTSansText>
      </View>
    );
  }

  render() {
    const { navigation, userProfile, messages } = this.props;
    const { isModalVisible, modalViewType } = this.state;

    let modal;
    switch(modalViewType) {
      case 'conclude':
        modal = this._renderConcludeModal();
        break;

      case 'completed':
        modal = this._renderCompletedModal();
        if (isModalVisible) {
          this.modalTimeout = setTimeout(() => this.setState({ isModalVisible: false }), 2000);
        } else {
          clearTimeout(this.modalTimeout);
        }
        break;
    }

    return (
      <View style={styles.container}>
        <MeetOverModal
          isVisible={isModalVisible}
          onBackdropPress={() => this.setState({ isModalVisible: false })}
        >
          {modal}
        </MeetOverModal>
        <GiftedChat
          messages={messages}
          onSend={messages => this._handleSend(messages)}
          messageIdGenerator={() => null}
          user={{
            _id: userProfile.id,
            name: userProfile.formattedName,
            avatar: userProfile.pictureUrl
          }}
          loadEarlier={messages && messages[messages.length - 1]._id !== 0}
          onLoadEarlier={() => this._handleLoadEarlier()}
          onPressAvatar={() => this._handlePressAvatar()}
          renderLoading={() => <Spinner color={Colors.tintColor} />}
          renderBubble={this.renderBubble}
          inverted={true}
        />
      </View>
    );
  }

  componentDidMount() {
    const {
      messages,
      navigation,
      registerFetchNewMessageAsync,
      fetchEarlierMessagesAsync
    } = this.props;

    navigation.setParams({ _handleConclude: () => this.setState({ isModalVisible: true }) });

    if (!messages) {
      fetchEarlierMessagesAsync(navigation.state.params._id, chatMessagesToLoad);
      registerFetchNewMessageAsync(navigation.state.params._id);
    }
  }
};

const mapStateToProps = (state, ownProps) => ({
  userProfile: state.userProfile,
  messages: state.chat.messageThreads[ownProps.navigation.state.params._id]
});

const mapDispatchToProps = {
  registerFetchNewMessageAsync,
  fetchEarlierMessagesAsync
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ChatScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  headerButton: {
    alignSelf: 'center'
  },
  headerText: {
    color: '#fff'
  },
  modalButton: {
    marginTop: 20,
    backgroundColor: Colors.tintColor,
    alignSelf: 'center'
  },
  modalHeaderText: {
    alignSelf: 'center'
  },
  modalText: {
    fontSize: 18
  }
});

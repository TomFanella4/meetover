import React from 'react';
import { StyleSheet, Platform } from 'react-native';
import {
  View,
  Button,
  Icon,
  Container,
  Form,
  Textarea,
  Spinner
} from 'native-base';
import Modal from 'react-native-modal';
import { connect } from 'react-redux';

import MeetOverModal from '../components/MeetOverModal';
import Profile from '../components/Profile';
import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';
import { separator, serverURI } from '../constants/Common';
import { requestScreenStrings } from '../constants/Strings';
import Layout from '../constants/Layout';
import { StyledToast } from '../helpers';

class RequestScreen extends React.Component {
  static navigationOptions = ({ navigation }) => ({
    title: `${navigation.state.params.profile.formattedName}'s Profile`,
  });

  state = {
    isModalVisible: false,
    modalText: '',
    modalViewType: 'request',
    thread: undefined
  };

  _loadThread(prevThreadList) {
    const { threadList, navigation, signedInProfile } = this.props;
    const id = navigation.state.params.profile.id;
    const signedInId = signedInProfile.id;
    let threadId;

    // Thread List has not been fetched from firebase yet
    if (threadList === null) {
      return;
    }

    // Get threadId
    if (signedInId < id) {
      threadId = signedInId + separator + id;
    } else {
      threadId = id + separator + signedInId;
    }

    // The thread has not changed
    if (prevThreadList && prevThreadList[threadId] === threadList[threadId]) {
      return;
    }

    // Update the current thread
    this.setState({ thread: threadList[threadId] });
  }

  componentWillMount() {
    this._blurSub = this.props.navigation.addListener('didBlur', () => this._onBlur());
  }

  componentDidMount() {
    this._loadThread();
  }

  componentDidUpdate(prevProps) {
    this._loadThread(prevProps.threadList);
  }

  _onBlur() {
    this.setState({ meetOverLoading: false });
  }

  _renderRequestModal() {
    const { modalText } = this.state;
    return (
      <View>
        <PTSansText style={styles.modalText}>Message</PTSansText>
        <Form style={styles.form}>
          <Textarea
            value={modalText}
            onChangeText={modalText => this.setState({ modalText })}
            rowSpan={5}
            bordered
            placeholder={requestScreenStrings.initialMessage}
          />
        </Form>
        <Button
          style={styles.confirmRequestButton}
          onPress={() => this._initiateMeetover()}
        >
          <PTSansText style={styles.request}>Request</PTSansText>
        </Button>
      </View>
    );
  }

  _renderLoadingModal() {
    return (
      <View>
        <PTSansText style={styles.modalText}>Sending MeetOver Request</PTSansText>
        <Spinner color={Colors.tintColor} />
      </View>
    );
  }

  _renderSuccessModal() {
    return (
      <View style={styles.modalSuccessContent}>
        <PTSansText style={styles.modalText}>Request Sent</PTSansText>
        <Icon style={styles.successIcon} name='checkmark-circle' />
      </View>
    );
  }

  async _initiateMeetover() {
    const { navigation, signedInProfile } = this.props;
    const { thread, modalText } = this.state;
    const id = navigation.state.params.profile.id;
    const signedInId = signedInProfile.id;
    const accessToken = signedInProfile.token.access_token;

    this.setState({ modalViewType: 'loading' });

    if (thread === undefined || (thread.status === 'declined' && thread.origin === 'receiver')) {
      const initialMessage = modalText !== '' ? modalText : requestScreenStrings.initialMessage;
      const uri = `${serverURI}/meetover/${id}`;
      const init = {
        method: 'POST',
        body: JSON.stringify({ initialMessage }),
        headers: new Headers({
          'Token': accessToken,
          'Identity': signedInId
        })
      };

      const response = await fetch(uri, init)
        .catch(err => console.log(err));

      if (response.status !== 200) {
        console.log('Could not initiate meetover');
        console.log(response);

        StyledToast({
          text: 'Could not initate MeetOver',
          buttonText: 'Okay',
          type: 'danger',
          duration: 3000,
        });

        return;
      }

      this.setState({ modalViewType: 'success' });

    } else {
      StyledToast({
        text: 'Could not initate MeetOver',
        buttonText: 'Okay',
        type: 'danger',
        duration: 3000,
      });
    }
  };

  _handleRequestButtonPress() {
    const { navigation } = this.props;
    const { thread } = this.state;

    if (thread) {
      switch(thread.status) {
        case 'pending':
          if (thread.origin === 'receiver') {
            navigation.navigate('ConfirmScreen', {
              _id: thread._id,
              profile: thread.profile
            });
            return;
          }
          break;

        case 'accepted':
          navigation.navigate('ChatScreen', {
            _id: thread._id,
            profile: thread.profile
          });
          return;

        case 'declined':
          if (thread.origin === 'sender') {
            return;
          }
          break;
      }
    }

    this.setState({ isModalVisible: true });
  }

  render() {
    const { profile, matchStatus } = this.props.navigation.state.params;
    const { isModalVisible, modalViewType, thread } = this.state;

    let modalView;

    switch(modalViewType) {
      case 'request':
        modalView = this._renderRequestModal();
        break;

      case 'loading':
        modalView = this._renderLoadingModal();
        break;

      case 'success':
        modalView = this._renderSuccessModal();
        isModalVisible && setTimeout(() => this.setState({ isModalVisible: false }), 2000);
        break;
    }

    let buttonText = 'Request MeetOver';
    let buttonDisabled;
    if (thread) {
      switch(thread.status) {
        case 'pending':
          if (thread.origin === 'receiver') {
            buttonText = 'View Request';
          } else {
            buttonText = 'Request Sent';
            buttonDisabled = true;
          }
          break;

        case 'accepted':
          buttonText = 'Open Chat';
          break;

        case 'declined':
          if (thread.origin === 'sender') {
            buttonText = 'Unavailable';
            buttonDisabled = true;
          }
          break;
      }
    }

    return (
      <Container style={styles.container}>
        <MeetOverModal
          isVisible={isModalVisible}
          onBackdropPress={() => this.setState({ isModalVisible: false })}
          onModalHide={() => this.setState({ modalViewType: 'request' })}
        >
          {modalView}
        </MeetOverModal>
        <Profile profile={profile} matchStatus={matchStatus} />
        <Button
          iconLeft
          full
          style={!buttonDisabled ? styles.requestButton : null}
          disabled={buttonDisabled}
          onPress={() => this._handleRequestButtonPress()}
        >
          <Icon name='chatboxes' />
          <PTSansText style={styles.request}>{buttonText}</PTSansText>
        </Button>
      </Container>
    );
  }

  componentWillUnmount() {
    this._blurSub.remove();
  }
};

const mapStateToProps = state => ({
  signedInProfile: state.userProfile,
  threadList: state.chat.threadList
});

export default connect(
  mapStateToProps
)(RequestScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  loadingView: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center'
  },
  modalSuccessContent: {
    flexDirection: 'row',
    justifyContent: 'center',
  },
  modalText: {
    paddingRight: 10,
    alignSelf: 'center'
  },
  successIcon: {
    color: Colors.success
  },
  form: {
    paddingTop: 5,
    paddingBottom: 10
  },
  requestButton: {
    backgroundColor: Colors.tintColor
  },
  confirmRequestButton: {
    backgroundColor: Colors.tintColor,
    alignSelf: 'center'
  },
  request: {
    fontSize: 18
  }
});

import React from 'react';
import { StyleSheet, TouchableWithoutFeedback, Keyboard } from 'react-native';
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
import { find } from 'lodash';

import Profile from '../components/Profile';
import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';
import { separator, serverURI } from '../constants/Common';
import { StyledToast } from '../helpers';

class RequestScreen extends React.Component {
  static navigationOptions = ({ navigation }) => ({
    title: `${navigation.state.params.profile.formattedName}'s Profile`,
  });

  state = {
    isModalVisible: false,
    isKeyboardVisible: false,
    modalText: '',
    modalViewType: 'request',
    thread: undefined
  };

  _loadThread() {
    const { threadList, navigation, signedInProfile } = this.props;
    const id = navigation.state.params.profile.id;
    const signedInId = signedInProfile.id;
    let threadId;

    if (signedInId < id) {
      threadId = signedInId + separator + id;
    } else {
      threadId = id + separator + signedInId;
    }

    const thread = find(threadList, { '_id': threadId });
    this.setState({ thread });
  }

  componentWillMount() {
    this.keyboardDidShowListener = Keyboard.addListener('keyboardDidShow', () => this._keyboardDidShow());
    this.keyboardDidHideListener = Keyboard.addListener('keyboardDidHide', () => this._keyboardDidHide());
    this._blurSub = this.props.navigation.addListener('didBlur', () => this._onBlur());
  }

  componentDidMount() {
    if (this.props.threadList !== null) {
      // Thread list has been fetched from Firebase
      this._loadThread();
    }
  }

  componentDidUpdate(prevProps) {
    if ((prevProps.threadList === null && this.props.threadList !== null) ||
         prevProps.threadList.length !== this.props.threadList.length ) {
      // Thread list has been fetched from Firebase
      this._loadThread();
    }
  }

  _onBlur() {
    this.setState({ meetOverLoading: false });
  }

  _keyboardDidShow() {
    this.setState({ isKeyboardVisible: true });
  }

  _keyboardDidHide() {
    this.setState({ isKeyboardVisible: false });
  }

  _renderRequestModal() {
    const { modalText } = this.state;
    return (
      <View style={styles.modalContent}>
        <PTSansText style={styles.modalText}>Request Message</PTSansText>
        <Form style={styles.form}>
          <Textarea
            value={modalText}
            onChangeText={modalText => this.setState({ modalText })}
            rowSpan={5}
            bordered
            placeholder={'Hello, I\'d like to MeetOver!'}
          />
        </Form>
        <Button
          style={styles.confirmRequestButton}
          onPress={() => this._initiateMeetover()}
        >
          <PTSansText style={styles.request}>Request MeetOver</PTSansText>
        </Button>
      </View>
    );
  }

  _renderLoadingModal() {
    return (
      <View style={styles.modalContent}>
        <PTSansText style={styles.modalText}>Sending MeetOver Request...</PTSansText>
        <Spinner color={Colors.tintColor} />
      </View>
    );
  }

  _renderSuccessModal() {
    return (
      <View style={styles.modalContent}>
        <PTSansText style={styles.modalText}>Success!</PTSansText>
      </View>
    );
  }

  _handleBackdropPress() {
    const { isKeyboardVisible, modalViewType } = this.state;
    if (isKeyboardVisible) {
      Keyboard.dismiss();
    } else if (modalViewType === 'request') {
      this.setState({ isModalVisible: false });
    }
  }

  async _initiateMeetover() {
    const { navigation, signedInProfile } = this.props;
    const { thread } = this.state;
    const id = navigation.state.params.profile.id;
    const signedInId = signedInProfile.id;
    const accessToken = signedInProfile.token.access_token;

    this.setState({ modalViewType: 'loading' });

    if (thread === undefined || (thread.status === 'declined' && thread.origin === 'receiver')) {
      const uri = `${serverURI}/meetover/${id}`;
      const init = {
        method: 'POST',
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
    const { thread } = this.state;

    thread === undefined ?
      this.setState({ isModalVisible: true })
    :
      this.props.navigation.navigate('ChatScreen', {
        _id: thread._id,
        profile: thread.profile
      });
  }

  render() {
    const { profile } = this.props.navigation.state.params;
    const { isModalVisible, modalViewType, thread } = this.state;
    const buttonDisabled = thread !== undefined && thread.status === 'pending';

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
        setTimeout(() => this.setState({ isModalVisible: false, }), 2000);
        break;
    }

    let buttonText;
    if (thread === undefined) {
      buttonText = 'Request MeetOver'
    } else if (thread.status === 'pending') {
      buttonText = 'Awaiting Response'
    } else if (thread.status === 'accepted') {
      buttonText = 'Open Chat'
    }

    return (
      <Container style={styles.container}>
        <Modal
          isVisible={isModalVisible}
          onBackdropPress={() => this._handleBackdropPress()}
          onModalHide={() => this.setState({ modalViewType: 'request' })}
        >
          <TouchableWithoutFeedback onPress={Keyboard.dismiss}>
            {modalView}
          </TouchableWithoutFeedback>
        </Modal>
        <Profile profile={profile} />
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
    this.keyboardDidShowListener.remove();
    this.keyboardDidHideListener.remove();
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
  modalContent: {
    backgroundColor: 'white',
    padding: 22,
    borderRadius: 4,
    borderColor: 'rgba(0, 0, 0, 0.1)',
  },
  modalText: {
    alignSelf: 'center'
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

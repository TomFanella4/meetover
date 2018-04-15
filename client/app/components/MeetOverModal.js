import React from 'react';
import {
  StyleSheet,
  TouchableWithoutFeedback,
  Keyboard,
} from 'react-native';
import { View } from 'native-base'
import Modal from 'react-native-modal';
import { connect } from 'react-redux';

class MeetOverModal extends React.Component {
  state = {
    isKeyboardVisible: false
  };

  componentWillMount() {
    this.keyboardDidShowListener = Keyboard.addListener('keyboardDidShow', () => this._keyboardDidShow());
    this.keyboardDidHideListener = Keyboard.addListener('keyboardDidHide', () => this._keyboardDidHide());
  }

  _keyboardDidShow() {
    this.setState({ isKeyboardVisible: true });
  }

  _keyboardDidHide() {
    this.setState({ isKeyboardVisible: false });
  }

  _handleBackdropPress() {
    if (this.state.isKeyboardVisible) {
      Keyboard.dismiss();
    } else {
      this.props.onBackdropPress();
    }
  }

  render() {
    return (
      <Modal
        isVisible={this.props.isVisible}
        onBackdropPress={() => this._handleBackdropPress()}
        onModalHide={this.props.onModalHide}
        avoidKeyboard={true}
      >
        <TouchableWithoutFeedback onPress={Keyboard.dismiss}>
          <View style={styles.modalContent}>
            {this.props.children}
          </View>
        </TouchableWithoutFeedback>
      </Modal>
    );
  }

  componentWillUnmount() {
    this.keyboardDidShowListener.remove();
    this.keyboardDidHideListener.remove();
  }
};

export default MeetOverModal;

const styles = StyleSheet.create({
  modalContent: {
    backgroundColor: 'white',
    padding: 22,
    borderRadius: 4,
    borderColor: 'rgba(0, 0, 0, 0.1)',
  }
});

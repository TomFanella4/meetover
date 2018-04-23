import React from 'react';
import { Platform, StyleSheet, Keyboard } from 'react-native';
import {
  View,
  Switch,
  Form,
  Textarea,
  Item,
  Label,
  Input,
  Button,
  Icon
} from 'native-base';
import { connect } from 'react-redux';
import DateTimePicker from 'react-native-modal-datetime-picker';
import moment from 'moment';

import MeetOverModal from './MeetOverModal';
import { modifyFirebaseUserState } from '../firebase';
import { modifyProfile } from '../actions/userActions';
import { PTSansText } from './StyledText';
import Colors from '../constants/Colors';
import { isSearchingBarStrings } from '../constants/Strings';

class IsSearchingBar extends React.Component {

  state = {
    isModalVisible: this.props.isModalVisible,
    isDateTimePickerVisible: false,
    greeting: '',
    timeAvailable: moment(),
    origin: '',
    destination: ''
  }

  _showDateTimePicker = () => {
    Keyboard.dismiss();
    this.setState({ isDateTimePickerVisible: true });
  };

  _hideDateTimePicker = () => this.setState({ isDateTimePickerVisible: false });

  _handleDatePicked = (date) => {
    let timeAvailable = moment(date.getTime());
    if (timeAvailable < moment()) {
      timeAvailable.add(1, 'd');
    }
    this.setState({ timeAvailable });
    this._hideDateTimePicker();
  };

  _renderModal() {
    const { greeting, timeAvailable, origin, destination } = this.state;
    return (
      <View>
        <PTSansText style={styles.modalText}>Start Searching</PTSansText>
        <Textarea
          value={greeting}
          onChangeText={greeting => this.setState({ greeting })}
          rowSpan={3}
          bordered
          placeholder={isSearchingBarStrings.initialGreeting}
        />
        <Form style={styles.form}>
          <Item onPress={this._showDateTimePicker}>
            <Label style={styles.formItem}>End Time: </Label>
            <Icon name='calendar' />
            <PTSansText>{timeAvailable.format('MMM Do h:mm a')}</PTSansText>
          </Item>
          <Item>
            <Label style={styles.formItem}>Origin:</Label>
            <Input
              style={styles.formItem}
              placeholder={'MIA'}
              value={origin}
              onChangeText={origin => this.setState({ origin })}
            />
          </Item>
          <Item>
            <Label style={styles.formItem}>Destination:</Label>
            <Input
              style={styles.formItem}
              placeholder={'JFK'}
              value={destination}
              onChangeText={destination => this.setState({ destination })}
            />
          </Item>
        </Form>
        <Button
          style={styles.modalButton}
          onPress={() => {
            Keyboard.dismiss();
            this.setState({ isModalVisible: false });
            const matchState = {
              isSearching: true,
              isMatched: false,
              greeting: greeting !== '' ? greeting : isSearchingBarStrings.initialGreeting,
              timeAvailable: timeAvailable.valueOf(),
              origin,
              destination
            }
            modifyFirebaseUserState('matchStatus', matchState);
          }}
        >
          <PTSansText style={styles.modalButtonText}>Start MeetOver</PTSansText>
        </Button>
      </View>
    );
  }

  render() {
    const { userProfile, modifyProfile } = this.props;
    return (
      <View style={styles.container}>
        <MeetOverModal
          isVisible={this.state.isModalVisible}
          onBackdropPress={() => {
            this.setState({ isModalVisible: false });
            modifyProfile('isSearching', false);
          }}
        >
          {this._renderModal()}
          <DateTimePicker
            isVisible={this.state.isDateTimePickerVisible}
            onConfirm={this._handleDatePicked}
            onCancel={this._hideDateTimePicker}
            mode={'time'}
            is24Hour={false}
          />
        </MeetOverModal>
        <PTSansText style={styles.isSearchingText}>
          {
            userProfile.isSearching ?
              isSearchingBarStrings.searchingText
            :
              isSearchingBarStrings.startText
          }
        </PTSansText>
        <Switch
          style={styles.isSearchingSwitch}
          value={userProfile.isSearching}
          onValueChange={value => {
            modifyProfile('isSearching', value);
            value && this.setState({ isModalVisible: true });
            !value && modifyFirebaseUserState('matchStatus/isSearching', false);
          }}
          onTintColor={Colors.tintColor}
          thumbTintColor={Platform.OS === 'android' ? 'white' : null}
        />
      </View>
    );
  }
};

const mapStateToProps = state => ({
  userProfile: state.userProfile
});

const mapDispatchToProps = {
  modifyProfile
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(IsSearchingBar);

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    padding: 10,
    backgroundColor: '#fff',
    borderBottomColor: '#bbb',
    borderBottomWidth: StyleSheet.hairlineWidth,
  },
  isSearchingText: {
    flex: 8,
  },
  isSearchingSwitch: {
    flex: 2,
  },
  modalText: {
    alignSelf: 'center'
  },
  form: {
    paddingTop: 5,
    paddingBottom: 10
  },
  modalButton: {
    backgroundColor: Colors.tintColor,
    alignSelf: 'center'
  },
  modalButtonText: {
    fontSize: 18
  },
  formItem: {
    fontFamily: 'pt-sans'
  }
});

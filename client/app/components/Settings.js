import React from 'react';
import { Platform, StyleSheet } from 'react-native';
import { connect } from 'react-redux';
import {
  Item,
  Label,
  Input,
  Switch,
  Form,
} from 'native-base';

import Colors from '../constants/Colors';
import { settingsStrings } from '../constants/Strings';
import { PTSansText } from '../components/StyledText';
import { modifyProfile } from '../actions/userActions';

const Settings = ({ userProfile, formOptions, modifyProfile, onProfileModified }) => {

  // TODO: Add check for empty user profile fields
  const formItems = formOptions.map(option => (
    <Item key={option.item} floatingLabel>
      <Label
        style={{
          color: !userProfile.isAuthenticated ? 'white' : null,
          fontFamily: 'pt-sans'
        }}
      >
        {option.label}
      </Label>
      <Input
        style={{
          color: !userProfile.isAuthenticated ? 'white' : null,
          fontFamily: 'pt-sans'
        }}
        value={userProfile[option.item] || ''}
        onChangeText={text => modifyProfile(option.item, text)}
        onBlur={() => (
          onProfileModified &&
          onProfileModified(
            option.item,
            userProfile[option.item]
          )
        )}
      />
    </Item>
  ));

  formItems.push(
    <Item style={styles.shareLocationView} key={'shareLocation'} last>
      <PTSansText style={
        userProfile.isAuthenticated ?
          styles.shareLocationText
        :
          styles.shareLocationWhiteText
      }>
        {settingsStrings.permission}
      </PTSansText>
      <Switch
        style={styles.shareLocationSwitch}
        value={userProfile.shareLocation}
        onValueChange={value => {
          modifyProfile('shareLocation', value);
          onProfileModified && onProfileModified('shareLocation', value);
        }}
        onTintColor={
          userProfile.isAuthenticated ?
            Colors.tintColor
          :
            'white'
        }
        thumbTintColor={Platform.OS === 'android' ? 'white' : null}
      />
    </Item>
  );

  return (
    <Form>
      {formItems}
    </Form>
  );
}

const mapStateToProps = state => ({
  userProfile: state.userProfile
});

const mapDispatchToProps = {
  modifyProfile
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Settings);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  shareLocationView: {
    flexDirection: 'row',
    paddingTop: 10,
    paddingBottom: 10
  },
  shareLocationText: {
    flex: 8,
    paddingRight: 5,
  },
  shareLocationWhiteText: {
    color: 'white',
    flex: 8,
    paddingRight: 5,
  },
  shareLocationSwitch: {
    alignSelf: 'center',
    flex: 2,
    paddingLeft: 5,
  },
});

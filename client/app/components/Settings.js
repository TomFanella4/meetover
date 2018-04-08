import React from 'react';
import { Platform, StyleSheet } from 'react-native';
import { connect } from 'react-redux';
import {
  View,
  Item,
  Label,
  Input,
  Button,
  Switch,
  Form,
  Thumbnail
} from 'native-base';

import Colors from '../constants/Colors';
import { settingsStrings } from '../constants/Strings';
import { PTSansText } from '../components/StyledText';
import { modifyProfile } from '../actions/userActions';

const Settings = ({ userProfile, formOptions, modifyProfile, onProfileModified, navigation }) => {
  const formItems = [];
  const formItemStyle = {
    color: !userProfile.isAuthenticated ? 'white' : null,
    fontFamily: 'pt-sans',
  };

  formItems.push(
    userProfile.pictureUrl !== '' ?
      <Thumbnail
        style={styles.picture}
        key='picture'
        source={{ uri: userProfile.pictureUrl }}
      />
    :
      <Thumbnail
        style={styles.picture}
        key='picture'
        source={require('../../assets/images/icon.png')}
      />
  );

  formItems.push(formOptions.map(option => (
    <Item key={option.item} stackedLabel>
      <Label style={formItemStyle}>
        {option.label}
      </Label>
      <Input
        style={formItemStyle}
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
  )));

  formItems.push(
    <Item style={styles.shareLocationView} key='shareLocation'>
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
  picture: {
    marginTop: 10,
    alignSelf: 'center'
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

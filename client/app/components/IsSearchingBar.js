import React from 'react';
import { Platform, StyleSheet } from 'react-native';
import { View, Switch } from 'native-base';
import { connect } from 'react-redux';

import { modifyUserProfile } from '../actions/index';
import { PTSansText } from './StyledText';
import Colors from '../constants/Colors';
import { isSearchingBarStrings } from '../constants/Strings';

const IsSearchingBar = ({ userProfile, modifyUserProfile }) => (
  <View style={styles.container}>
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
      onValueChange={value => modifyUserProfile({ isSearching: value })}
      onTintColor={Colors.tintColor}
      thumbTintColor={Platform.OS === 'android' ? 'white' : null}
    />
  </View>
);

const mapStateToProps = state => ({
  userProfile: state.userProfile
});

const mapDispatchToProps = {
  modifyUserProfile
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
});

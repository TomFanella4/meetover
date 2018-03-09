import React from 'react';
import { Platform, StyleSheet } from 'react-native';
import { connect } from 'react-redux';
import Expo from 'expo';
import {
  Container,
  Content,
  Button,
} from 'native-base';

import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';
import Settings from '../components/Settings';
import { StyledToast } from '../helpers';
import { deleteProfileAndLogoutAsync } from '../actions';
import { modifyFirebaseUserProfile } from '../firebase';

class CreateProfileScreen extends React.Component {
  static navigationOptions = {
    title: 'Settings',
  };

  formOptions = [
    { label: 'Name', item: 'formattedName' },
    { label: 'Headline', item: 'headline' },
    { label: 'Position', item: 'position' },
    { label: 'Company', item: 'company' },
    { label: 'Industry', item: 'industry' },
  ];

  render() {
    return (
      <Container style={styles.container}>
        <Content>
          <Settings
            formOptions={this.formOptions}
            onProfileModified={this._handleProfileModification.bind(this)}
          />
          <Button
            onPress={() => this._handleSignOutButtonPress()}
            style={styles.signOutButton}
          >
            <PTSansText>
              Sign Out
            </PTSansText>
          </Button>
        </Content>
      </Container>
    );
  }

  _handleProfileModification(key, value) {
    const { userProfile } = this.props;
    Promise.all([
      Expo.SecureStore.setItemAsync(
        'userProfile',
        JSON.stringify(userProfile)
      ),
      modifyFirebaseUserProfile(key, value),
    ])
    .then(() => StyledToast({
      text: 'Saved Settings',
      buttonText: 'Okay',
    }))
    .catch(() => StyledToast({
      text: 'Failed to save settings',
      buttonText: 'Okay',
      type: 'danger',
    }));
  }

  _handleSignOutButtonPress() {
    const { deleteProfileAndLogoutAsync, userProfile, navigation } = this.props;
    deleteProfileAndLogoutAsync();
    navigation.navigate('Login');
    StyledToast({ text: `${userProfile.firstName} Signed Out` });
  }
}

const mapStateToProps = state => ({
  userProfile: state.userProfile
});

const mapDispatchToProps = {
  deleteProfileAndLogoutAsync
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(CreateProfileScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  signOutButton: {
    alignSelf: 'center',
    marginTop: 10,
    marginBottom: 10,
    backgroundColor: Colors.tintColor,
  },
});

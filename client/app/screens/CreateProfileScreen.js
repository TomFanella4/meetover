import React from 'react';
import { Platform, StyleSheet } from 'react-native';
import { connect } from 'react-redux';
import {
  Container,
  View,
  Content,
  Button,
} from 'native-base';

import Colors from '../constants/Colors';
import { createProfileScreenStrings } from '../constants/Strings';
import { PTSansText } from '../components/StyledText';
import Settings from '../components/Settings';
import { StyledToast } from '../helpers';
import { saveProfileAndLoginAsync } from '../actions';

class CreateProfileScreen extends React.Component {
  static navigationOptions = {
    header: null
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
            <PTSansText style={styles.titleText}>
              {createProfileScreenStrings.title}
            </PTSansText>
            <PTSansText style={styles.titleSubtext}>
              {createProfileScreenStrings.subtitle}
            </PTSansText>
            <Settings formOptions={this.formOptions} />
        </Content>
        <View style={styles.finishView}>
          <Button
            style={styles.finishButton}
            onPress={() => this._handleFinishButtonPress()}
          >
            <PTSansText>Finish</PTSansText>
          </Button>
          <PTSansText style={styles.termsText}>
            {createProfileScreenStrings.notification}
          </PTSansText>
        </View>
      </Container>
    );
  }

  _handleFinishButtonPress() {
    const { saveProfileAndLoginAsync, userProfile } = this.props;
    let error = false;

    this.formOptions.forEach(formOption => {
      if (userProfile[formOption.item] === undefined ||
          userProfile[formOption.item] === '') {
        error = true;
      }
    });

    if (error) {
      StyledToast({
        text: 'Please Fill in All Fields',
        buttonText: 'Okay',
        type: 'danger',
        duration: 3000,
      });
      return;
    }

    saveProfileAndLoginAsync({ ...userProfile, isAuthenticated: true });
    StyledToast({ text: `${userProfile.firstName} Signed In` });
  }
}

const mapStateToProps = state => ({
  userProfile: state.userProfile
});

const mapDispatchToProps = {
  saveProfileAndLoginAsync
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(CreateProfileScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: Colors.tintColor,
    paddingTop: Platform.OS === 'ios' ? 15 : 0,
  },
  finishView: {
    backgroundColor: '#fff',
    padding: 20,
  },
  titleText: {
    color: 'white',
    fontSize: 29,
    alignSelf: 'center',
    paddingTop: 15
  },
  titleSubtext: {
    color: 'white',
    fontSize: 20,
    alignSelf: 'center',
  },
  termsText: {
    color: Colors.tintColor,
    fontSize: 10,
    alignSelf: 'center',
  },
  finishButton: {
    backgroundColor: Colors.tintColor,
    alignSelf: 'center',
  },
});

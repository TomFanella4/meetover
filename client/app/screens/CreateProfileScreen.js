import React from 'react';
import { Platform, StyleSheet } from 'react-native';
import {
  Container,
  View,
  Form,
  Item,
  Label,
  Input,
  Content,
  Button,
  Switch,
} from 'native-base';

import Colors from '../constants/Colors';
import { createProfileScreenStrings } from '../constants/Strings';
import { PTSansText } from '../components/StyledText';

import { connect } from 'react-redux';
import { saveProfileAndLoginAsync, modifyUserProfile } from '../actions';

class CreateProfileScreen extends React.Component {
  static navigationOptions = {
    header: null
  };

  userProfileFormOptions = [
    { label: 'Name', item: 'formattedName' },
    { label: 'Headline', item: 'headline' },
    { label: 'Position', item: 'position' },
    { label: 'Company', item: 'company' },
    { label: 'Industry', item: 'industry' },
  ];

  render() {
    const { userProfile } = this.props;

    const userProfileFormItems = this.userProfileFormOptions.map(option => (
      <Item key={option.item} floatingLabel last>
        <Label
          style={{color: 'white', fontFamily: 'pt-sans'}}
        >
          {option.label}
        </Label>
        <Input
          style={{color: 'white', fontFamily: 'pt-sans'}}
          value={userProfile[option.item] || ''}
          onChangeText={text => this._handleUserProfileModification(option.item, text)}
        />
      </Item>
    ));

    return (
      <Container style={styles.container}>
        <Content style={styles.contentTop}>
            <PTSansText style={styles.titleText}>
              {createProfileScreenStrings.title}
            </PTSansText>
            <PTSansText style={styles.titleSubtext}>
              {createProfileScreenStrings.subtitle}
            </PTSansText>
            <Form>
              {userProfileFormItems}
              <View style={styles.shareLocationView}>
                <PTSansText style={styles.shareLocationText}>
                  {createProfileScreenStrings.permission}
                </PTSansText>
                <Switch style={styles.shareLocationSwitch}
                  value={userProfile.shareLocation}
                  onValueChange={value => this._handleUserProfileModification('shareLocation', value)}
                  onTintColor='white'
                  thumbTintColor={Platform.OS === 'android' ? 'white' : null} />
              </View>
            </Form>
        </Content>
        <View style={styles.contentBottom}>
          <Button
            style={styles.button}
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

  _handleUserProfileModification(key, value) {
    const { modifyUserProfile, userProfile } = this.props;
    modifyUserProfile({ ...userProfile, [key]: value });
  }

  _handleFinishButtonPress() {
    const { saveProfileAndLoginAsync, userProfile } = this.props;
    saveProfileAndLoginAsync({ ...userProfile, isAuthenticated: true });
  }
}

const mapStateToProps = state => ({
  userProfile: state.userProfile
});

const mapDispatchToProps = {
  saveProfileAndLoginAsync,
  modifyUserProfile
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
  contentTop: {
    paddingLeft: 20,
    paddingRight: 20,
  },
  contentBottom: {
    backgroundColor: '#fff',
    padding: 20,
  },
  shareLocationView: {
    flexDirection: 'row',
    paddingTop: 15,
    paddingBottom: 15
  },
  shareLocationText: {
    color: 'white',
    flex: 8,
    paddingRight: 5,
  },
  shareLocationSwitch: {
    alignSelf: 'center',
    flex: 2,
    paddingLeft: 5,
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
  button: {
    backgroundColor: Colors.tintColor,
    alignSelf: 'center',
  },
});

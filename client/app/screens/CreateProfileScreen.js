import React from 'react';
import { StyleSheet } from 'react-native';
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
  Text
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
        <Label style={{color: 'white'}}>{option.label}</Label>
        <Input
          style={{color: 'white'}}
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
                <Text style={styles.shareLocationText}>
                  {createProfileScreenStrings.permission}
                </Text>
                <Switch style={styles.shareLocationSwitch}
                  value={userProfile.shareLocation}
                  onValueChange={value => this._handleUserProfileModification('shareLocation', value)}
                  onTintColor='white' />
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
    backgroundColor: '#fff',
  },
  contentTop: {
    backgroundColor: Colors.tintColor,
    padding: 20,
    paddingTop: 40,
  },
  contentBottom: {
    padding: 20,
  },
  shareLocationView: {
    flexDirection: 'row',
    paddingTop: 15,
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
    fontSize: 35,
    alignSelf: 'center',
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

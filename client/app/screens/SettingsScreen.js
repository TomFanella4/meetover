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
  Switch,
} from 'native-base';

import Colors from '../constants/Colors';
import { settingsScreenStrings } from '../constants/Strings';
import { PTSansText } from '../components/StyledText';

import { connect } from 'react-redux';
import { modifyUserProfile } from '../actions';

class CreateProfileScreen extends React.Component {
  static navigationOptions = {
    title: 'Settings',
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
          <Form>
            {userProfileFormItems}
            <View style={styles.shareLocationView}>
              <PTSansText style={styles.shareLocationText}>
                {settingsScreenStrings.permission}
              </PTSansText>
              <Switch style={styles.shareLocationSwitch}
                value={userProfile.shareLocation}
                onValueChange={value => this._handleUserProfileModification('shareLocation', value)}
                onTintColor='white'
                thumbTintColor={Platform.OS === 'android' ? 'white' : null}
              />
            </View>
          </Form>
        </Content>
      </Container>
    );
  }

  _handleUserProfileModification(key, value) {
    const { modifyUserProfile, userProfile } = this.props;
    modifyUserProfile({ ...userProfile, [key]: value });
  }
}

const mapStateToProps = state => ({
  userProfile: state.userProfile
});

const mapDispatchToProps = {
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
  },
  contentTop: {
    paddingLeft: 20,
    paddingRight: 20,
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
});

import React from 'react';
import { StyleSheet } from 'react-native';
import { Container } from 'native-base';

import Profile from '../components/Profile';

class ProfileScreen extends React.Component {
  static navigationOptions = ({ navigation }) => ({
    title: `${navigation.state.params.profile.formattedName}'s Profile`,
  });

  render() {
    const { profile } = this.props.navigation.state.params;
    return (
      <Container style={styles.container}>
        <Profile profile={profile} />
      </Container>
    );
  }
};

export default ProfileScreen;

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
});

import React from 'react';
import { View, Button, StyleSheet } from 'react-native';
import { NavigationActions } from 'react-navigation'

import { bindActionCreators } from 'redux'
import { connect } from 'react-redux';

import { authenticateUser } from '../actions';

class LoginScreen extends React.Component {

  componentDidUpdate() {
    if (this.props.authenticated) {
      const resetAction = NavigationActions.reset({
        index: 0,
        actions: [
          NavigationActions.navigate({routeName: 'Main'})
        ]
      });
      this.props.navigation.dispatch(resetAction);
    }
  }

  render() {
    return (
      <View style={styles.container}>
        <Button
          title='Sign in with LinkedIn'
          onPress={this.props.authenticateUser}
        />
      </View>
    )
  }
}

const mapStateToProps = state => ({
  authenticated: state.userProfileReducer.authenticated
});

const mapDispatchToProps = dispatch => (
  bindActionCreators({ authenticateUser }, dispatch)
);

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(LoginScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingTop: 15,
    backgroundColor: '#fff',
  }
});

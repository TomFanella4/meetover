import React from 'react';
import { Button, Container, Content, Text } from 'native-base';
import { connect } from 'react-redux';

import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';

class ProfileScreen extends React.Component {
  static navigationOptions = {
    title: 'Profile',
    headerStyle: {
      backgroundColor: Colors.tintColor,
    },
    headerTintColor: Colors.header,
  };

  render() {
    return (
      <Container>
        <Content>
          <PTSansText>suop</PTSansText>
          <Button onPress={() => this.props.navigation.goBack()} title="Back">
            <PTSansText>Back</PTSansText>
          </Button>
        </Content>
      </Container>
    );
  }
}

export default connect(
  null,
  null
)(ProfileScreen);

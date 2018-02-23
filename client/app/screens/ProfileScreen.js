import React from 'react';
import { Container, Content } from 'native-base';
import { connect } from 'react-redux';

import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';

class ProfileScreen extends React.Component {
  static navigationOptions = {
    headerStyle: {
      backgroundColor: Colors.tintColor,
    },
    headerTintColor: Colors.header,
  };

  render() {
    const { navigation } = this.props;

    return (
      <Container>
        <Content>
          <PTSansText>Test Profile</PTSansText>
        </Content>
      </Container>
    );
  }
}

export default connect(
  null,
  null
)(ProfileScreen);

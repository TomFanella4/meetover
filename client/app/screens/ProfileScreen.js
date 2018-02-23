import React from 'react';
import { Container, Content, Spinner } from 'native-base';
import { connect } from 'react-redux';

import { fetchProfileAsync } from '../actions';
import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';

class ProfileScreen extends React.Component {
  state = {
    loading: true
  };

  componentDidMount() {
    const { navigation } = this.props;

    this.props.fetchProfileAsync(navigation.state.params.userId)
      .then(() => this.setState({ loading: false }));
  }

  render() {
    const { navigation, profile } = this.props;

    let content;

    if (this.state.loading) {
      content = (<Spinner color={Colors.tintColor} />);
    } else {
      content = (
        <PTSansText>{profile.formattedName}</PTSansText>
      );
    }

    return (
      <Container>
        <Container>
          {content}
        </Container>
      </Container>
    );
  }
};

const mapStateToProps = state => ({
  profile: state.matchList.profile
});

const mapDispatchToProps = {
  fetchProfileAsync
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ProfileScreen);

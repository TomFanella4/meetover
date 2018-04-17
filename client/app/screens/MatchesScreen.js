import React from 'react';
import { RefreshControl, StyleSheet } from 'react-native';
import {
  Body,
  Container,
  Content,
  Left,
  List,
  ListItem,
  Spinner
} from 'native-base'
import { connect } from 'react-redux';

import { PTSansText } from '../components/StyledText';
import { ProfileImage } from '../components/ProfileImage';
import Colors from '../constants/Colors';
import IsSearchingBar from '../components/IsSearchingBar';
import { fetchMatchesAsync } from '../actions/matchesActions';

class MatchesScreen extends React.Component {
  static navigationOptions = {
    title: 'Matches',
  };

  state = {
    loading: true,
    refreshing: false
  };

  componentDidMount() {
    const { userId, accessToken } = this.props;
    this._onRefresh(userId, accessToken);
  }

  _onRefresh(userId, accessToken) {
    this.setState({ refreshing: true });
    this.props.fetchMatchesAsync(userId, accessToken)
      .then(() => this.setState({ loading: false, refreshing: false }));
  }

  _viewProfile(match) {
    this.props.navigation.navigate('RequestScreen', { profile: match });
  }

  render() {
    const { matches, userId, accessToken } = this.props;
    const list = matches.map(match => (
      <ListItem
        key={match.profile.id}
        onPress={() => this._viewProfile(match.profile)}
        avatar
      >
        <Left>
          <ProfileImage pictureUrl={match.profile.pictureUrl} />
        </Left>
        <Body>
          <PTSansText style={styles.name}>{match.profile.formattedName}</PTSansText>
          <PTSansText note>{match.profile.headline}</PTSansText>
        </Body>
      </ListItem>
    ));

    const refresh = (
      <RefreshControl
        onRefresh={() => this._onRefresh(userId, accessToken)}
        refreshing={this.state.refreshing}
      />
    );

    return (
      <Container style={styles.container}>
        <IsSearchingBar />
        {
          !this.state.loading ?
            <Content refreshControl={refresh}>
              <List>{list}</List>
            </Content>
          :
            <Spinner color={Colors.tintColor} />
        }
      </Container>
    );
  }
};

const mapStateToProps = state => ({
  userId: state.userProfile.id,
  accessToken: state.userProfile.token.access_token,
  matches: state.matchList.matches
});

const mapDispatchToProps = {
  fetchMatchesAsync
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(MatchesScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  name: {
    fontSize: 20
  }
});

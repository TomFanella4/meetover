import React from 'react';
import { RefreshControl, StyleSheet } from 'react-native';
import {
  Body,
  Container,
  Content,
  Left,
  List,
  ListItem,
  Thumbnail,
  Spinner
} from 'native-base'
import { connect } from 'react-redux';

import { PTSansText } from '../components/StyledText';
import Colors from '../constants/Colors';
import IsSearchingBar from '../components/IsSearchingBar';
import { fetchMatchesAsync } from '../actions/matchesActions';

class ListScreen extends React.Component {
  state = {
    loading: true,
    refreshing: false
  };

  static navigationOptions = {
    title: 'List',
  };

  componentDidMount() {
    this._onRefresh('userId');
  }

  _onRefresh(userId) {
    this.setState({ refreshing: true });
    this.props.fetchMatchesAsync(userId)
      .then(() => this.setState({ loading: false, refreshing: false }));
  }

  _viewProfile(userId, name) {
    const { navigation } = this.props;

    navigation.navigate('Profile', {
      userId,
      name,
    });
  }

  render() {
    const { matches } = this.props;
    const list = matches.map((match, index) =>
      <ListItem style={styles.container} key={index} onPress={() => this._viewProfile(match.id, match.formattedName)}>
        <Left style={styles.thumbnail}>
          <Thumbnail source={{ uri: match.pictureUrl }} />
        </Left>
        <Body>
          <PTSansText style={styles.name}>{match.formattedName}</PTSansText>
          <PTSansText style={styles.headline}>{match.headline}</PTSansText>
        </Body>
      </ListItem>
    );

    const refresh = (
      <RefreshControl
        onRefresh={() => this._onRefresh('userId')}
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
  matches: state.matchList.matches
});

const mapDispatchToProps = {
  fetchMatchesAsync
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ListScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  thumbnail: {
    flex: 0
  },
  name: {
    fontSize: 20
  },
  headline: {
    fontSize: 12
  }
});

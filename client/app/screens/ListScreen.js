import React from 'react';
import { RefreshControl, ScrollView, StyleSheet } from 'react-native';
import {
  Body,
  Container,
  Content,
  Left,
  List,
  ListItem,
  Thumbnail
} from 'native-base'
import { connect } from 'react-redux';

import { PTSansText } from '../components/StyledText';
import { fetchMatchesAsync } from '../actions';

class ListScreen extends React.Component {
  state = {
      refreshing: false
  };

  static navigationOptions = {
    title: 'List',
  };

  componentDidMount() {
    this.props.fetchMatchesAsync('userId');
  }

  _onRefresh(userId) {
    this.setState({ refreshing: true });
    this.props.fetchMatchesAsync(userId)
      .then(() => this.setState({ refreshing: false }));
  }

  render() {
    const { matches } = this.props;
    const list = matches.map(match =>
      <ListItem style={styles.container} key={match.id}>
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
        <Content refreshControl={refresh}>
          <List>{list}</List>
        </Content>
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

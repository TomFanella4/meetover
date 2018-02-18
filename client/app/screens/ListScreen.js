import React from 'react';
import { ScrollView, StyleSheet } from 'react-native';
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

import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';
import { fetchMatchesAsync } from '../actions';

class ListScreen extends React.Component {
  static navigationOptions = {
    title: 'List',
    headerStyle: {
      backgroundColor: Colors.tintColor,
    },
    headerTintColor: Colors.header,
  };

  componentDidMount() {
    this.props.fetchMatchesAsync('userId');
  }

  render() {
    const { matches } = this.props;
    const list = matches.map(match => {
      const position = match.positions.values[0];

      return (
        <ListItem style={styles.container} key={match.id}>
          <Left style={styles.thumbnail}>
            <Thumbnail source={{ uri: match.pictureUrl }} />
          </Left>
          <Body>
            <PTSansText style={styles.name}>{match.formattedName}</PTSansText>
            <PTSansText style={styles.title}>
              {position.title} at {position.company.name}
            </PTSansText>
          </Body>
        </ListItem>
      );
    });

    return (
      <Container style={styles.container}>
        <Content>
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
  title: {
    fontSize: 12
  }
});

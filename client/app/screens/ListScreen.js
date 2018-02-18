import React from 'react';
import { ScrollView, StyleSheet } from 'react-native';
import { Container, Content, List, ListItem } from 'native-base'
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
    const list = matches.map(match => <ListItem key={match.id}>
      <PTSansText>{match.formattedName}</PTSansText>
    </ListItem>);

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
  }
});

import React from 'react';
import { StyleSheet } from 'react-native';
import {
  Body,
  Button,
  Card,
  CardItem,
  Container,
  Content,
  Icon,
  Left,
  Spinner,
  Thumbnail,
} from 'native-base';
import { connect } from 'react-redux';

import { fetchProfileAsync } from '../actions/matchesActions';
import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';

class ProfileScreen extends React.Component {
  static navigationOptions = ({ navigation }) => ({
    title: `${navigation.state.params.name}'s Profile`,
  });

  state = {
    loading: true
  };

  componentDidMount() {
    const { fetchProfileAsync, navigation } = this.props;

    fetchProfileAsync(navigation.state.params.userId)
      .then(() => this.setState({ loading: false }));
  }

  _renderLoading() {
    return (
      <Container style={styles.container}>
        <Content>
          <Spinner color={Colors.tintColor} />
        </Content>
      </Container>
    );
  }

  render() {
    const { profile } = this.props;

    if (this.state.loading) {
      return this._renderLoading();
    } else {
      const positions = profile.positions.values.map((position, index) =>
        <CardItem style={styles.listItem} key={index}>
          <PTSansText style={styles.jobTitle}>
            {position.title} at {position.company.name}
          </PTSansText>
          <PTSansText>{position.summary}</PTSansText>
        </CardItem>
      );

      return (
        <Container style={styles.container}>
          <Content style={styles.container}>
            <Left style={styles.thumbnail}>
              <Thumbnail source={{ uri: profile.pictureUrl }} />
            </Left>
            <Body>
              <PTSansText style={styles.name}>{profile.formattedName}</PTSansText>
              <PTSansText>{profile.headline}</PTSansText>
              <PTSansText>
                <Icon name='pin' style={styles.location} /> {profile.location.name}
              </PTSansText>
            </Body>
            <Card>
              <CardItem header>
                <PTSansText style={styles.subtitle}>Summary</PTSansText>
              </CardItem>
              <CardItem><PTSansText>{profile.summary}</PTSansText></CardItem>
            </Card>
            <Card>
              <CardItem header>
                <PTSansText style={styles.subtitle}>Positions</PTSansText>
              </CardItem>
              {positions}
            </Card>
          </Content>
          <Button iconLeft full style={styles.chatButton}>
            <Icon name='chatboxes' />
            <PTSansText style={styles.request}>Request MeetOver</PTSansText>
          </Button>
        </Container>
      );
    }
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

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  listItem: {
    flexDirection: 'column',
  },
  location: {
    fontSize: 18,
    paddingRight: 2,
  },
  thumbnail: {
    paddingTop: 5,
    flex: 0
  },
  name: {
    fontSize: 26
  },
  chatButton: {
    backgroundColor: Colors.tintColor
  },
  request: {
    fontSize: 18
  },
  jobTitle: {
    alignSelf: 'flex-start',
    fontSize: 20
  },
  subtitle: {
    alignSelf: 'flex-start',
    fontSize: 22
  },
});

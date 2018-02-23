import React from 'react';
import { StyleSheet } from 'react-native';
import {
  Body,
  Card,
  CardItem,
  Container,
  Content,
  Left,
  Spinner,
  Thumbnail,
} from 'native-base';
import { connect } from 'react-redux';
import { Ionicons } from '@expo/vector-icons';

import { fetchProfileAsync } from '../actions';
import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';

class ProfileScreen extends React.Component {
  state = {
    loading: true
  };

  componentDidMount() {
    const { fetchProfileAsync, navigation } = this.props;

    fetchProfileAsync(navigation.state.params.userId)
      .then(() => this.setState({ loading: false }));
  }

  render() {
    const { navigation, profile } = this.props;

    let content;

    if (this.state.loading) {
      content = (
        <Content>
          <Spinner color={Colors.tintColor} />
        </Content>
      );
    } else {
      const positions = profile.positions.values.map((position, index) =>
        <CardItem style={styles.listItem} key={index}>
          <PTSansText style={styles.jobTitle}>
            {position.title} at {position.company.name}
          </PTSansText>
          <PTSansText>{position.summary}</PTSansText>
        </CardItem>
      );

      content = (
        <Content style={styles.container}>
          <Left style={styles.thumbnail}>
            <Thumbnail source={{ uri: profile.pictureUrl }} />
          </Left>
          <Body>
            <PTSansText style={styles.name}>{profile.formattedName}</PTSansText>
            <PTSansText>{profile.headline}</PTSansText>
            <PTSansText>
              <Ionicons name='md-pin' style={styles.location} /> {profile.location.name}
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
      );
    }

    return (
      <Container style={styles.container}>
        {content}
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

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  listItem: {
    flexDirection: 'column',
  },
  location: {
    paddingRight: 2,
  },
  thumbnail: {
    paddingTop: 5,
    flex: 0
  },
  name: {
    fontSize: 26
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

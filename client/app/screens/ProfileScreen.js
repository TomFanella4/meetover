import React from 'react';
import { StyleSheet } from 'react-native';
import {
  View,
  Body,
  Button,
  Card,
  CardItem,
  Container,
  Content,
  Icon,
  Left,
  Thumbnail,
  Spinner
} from 'native-base';
import { connect } from 'react-redux';
import { find } from 'lodash';

import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText';
import { separator, serverURI } from '../constants/Common';
import { StyledToast } from '../helpers';

class ProfileScreen extends React.Component {
  static navigationOptions = ({ navigation }) => ({
    title: `${navigation.state.params.profile.formattedName}'s Profile`,
  });

  state = {
    buttonDisabled: true,
    meetOverLoading: false
  };

  componentDidMount() {
    const { threadList } = this.props;

    if (threadList !== null) {
      // Thread list has been fetched from Firebase
      this.setState({ buttonDisabled: false });
    }
  }

  componentDidUpdate(prevProps) {
    const { threadList } = this.props;
    const prevThreadList = prevProps.threadList;

    if (prevThreadList === null && threadList !== null) {
      // Thread list has been fetched from Firebase
      this.setState({ buttonDisabled: false });
    }
  }

  _renderLoading() {
    // TODO: Add react navigation listener
    return (
      <Container style={styles.container}>
        <View style={styles.loadingView}>
          <PTSansText>Sending MeetOver Request...</PTSansText>
          <Spinner color={Colors.tintColor} />
        </View>
      </Container>
    );
  }

  async _initiateMeetover() {
    const { navigation, signedInProfile, threadList } = this.props;
    const { formattedName, id } = navigation.state.params.profile;
    const signedInId = signedInProfile.id;
    const accessToken = signedInProfile.token.access_token;
    let threadId;

    this.setState({ buttonDisabled: true, meetOverLoading: true });

    if (signedInId < id) {
      threadId = signedInId + separator + id;
    } else {
      threadId = id + separator + signedInId;
    }

    const exists = (find(threadList, { '_id': threadId }) !== undefined);

    if (!exists) {
      const uri = `${serverURI}/meetover/${id}`;
      const init = {
        method: 'POST',
        headers: new Headers({
          'Token': accessToken,
          'Identity': signedInId
        })
      };

      const response = await fetch(uri, init)
        .catch(err => console.log(err));

      if (response.status !== 200) {
        console.log('Could not initiate meetover');
        console.log(response);

        StyledToast({
          text: 'Could not initate MeetOver',
          buttonText: 'Okay',
          type: 'danger',
          duration: 3000,
        });

        return;
      }
    }

    this.setState({ buttonDisabled: false, meetOverLoading: false });
    navigation.navigate('ChatScreen', { _id: threadId, name: formattedName });
  };

  render() {
    const { profile } = this.props.navigation.state.params;
    const { buttonDisabled, meetOverLoading } = this.state;

    const positions = profile.positions.values.map((position, index) =>
      <CardItem style={styles.listItem} key={index}>
        <PTSansText style={styles.jobTitle}>
          {position.title} at {position.company.name}
        </PTSansText>
        <PTSansText>{position.summary}</PTSansText>
      </CardItem>
    );

    if(meetOverLoading) {
      return this._renderLoading();
    } else {
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
          <Button
            iconLeft
            full
            style={!buttonDisabled ? styles.chatButton : null}
            disabled={buttonDisabled}
            onPress={() => this._initiateMeetover()}
          >
            <Icon name='chatboxes' />
            <PTSansText style={styles.request}>Request MeetOver</PTSansText>
          </Button>
        </Container>
      );
    }
  }
};

const mapStateToProps = state => ({
  signedInProfile: state.userProfile,
  threadList: state.chat.threadList
});

export default connect(
  mapStateToProps
)(ProfileScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  loadingView: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center'
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

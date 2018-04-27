import React from 'react';
import { StyleSheet } from 'react-native';
import {
  Body,
  Card,
  CardItem,
  Content,
  Icon,
  Left
} from 'native-base';
import moment from 'moment';

import { PTSansText } from '../components/StyledText';
import { ProfileImage } from '../components/ProfileImage';

const Profile = ({ profile }) => {
  // // TODO cleanup data flow
  const { isSearching, isMatched, greeting, timeAvailable, origin, destination } = profile;
  const showMatchCard = true;

  const matchCard = showMatchCard && (
    <Card>
      <CardItem header>
        <PTSansText style={styles.subtitle}>Status</PTSansText>
      </CardItem>
      {
        greeting ?
          <CardItem><PTSansText>{greeting}</PTSansText></CardItem>
        :
          null
      }
      {
        (origin && destination) ?
          <CardItem style={styles.travelBox}>
            <PTSansText style={styles.travel}>{origin}</PTSansText>
            <Icon name='plane' style={styles.plane} />
            <PTSansText style={styles.travel}>{destination}</PTSansText>
          </CardItem>
        :
          null
      }
      {
        timeAvailable ?
          <CardItem>
            <PTSansText>Available until {moment(timeAvailable).format('MMMM Do, h:mm a')}.</PTSansText>
          </CardItem>
        :
          null
      }
    </Card>
  );

  const summaryCard = profile.summary !== '' && (
    <Card>
      <CardItem header>
        <PTSansText style={styles.subtitle}>Summary</PTSansText>
      </CardItem>
      <CardItem><PTSansText>{profile.summary}</PTSansText></CardItem>
    </Card>
  );

  const positionsCard = profile.positions._total > 0 && (
    <Card>
      <CardItem header>
        <PTSansText style={styles.subtitle}>Positions</PTSansText>
      </CardItem>
      {
        profile.positions.values.map(position => (
          <CardItem style={styles.listItem} key={position.id}>
            <PTSansText style={styles.jobTitle}>
              {position.title} at {position.company.name}
            </PTSansText>
            <PTSansText>{position.summary}</PTSansText>
          </CardItem>
        ))
      }
    </Card>
  );

  return (
    <Content>
      <Left style={styles.thumbnail}>
        <ProfileImage pictureUrl={profile.pictureUrl} />
      </Left>
      <Body>
        <PTSansText style={styles.name}>{profile.formattedName}</PTSansText>
        <PTSansText>{profile.headline}</PTSansText>
        <PTSansText>
          <Icon name='pin' style={styles.location} /> {profile.location.name}
        </PTSansText>
      </Body>
      {matchCard}
      {summaryCard}
      {positionsCard}
    </Content>
  );
};

export default Profile;

const styles = StyleSheet.create({
  listItem: {
    flexDirection: 'column',
  },
  location: {
    fontSize: 18,
    paddingRight: 2,
  },
  travelBox: {
    justifyContent: 'center',
  },
  travel: {
    fontSize: 20,
  },
  plane: {
    fontSize: 20,
    paddingRight: 10,
    paddingLeft: 10,
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
  }
});

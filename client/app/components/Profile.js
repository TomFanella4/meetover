import React from 'react';
import { StyleSheet } from 'react-native';
import {
  Body,
  Card,
  CardItem,
  Content,
  Icon,
  Left,
  Thumbnail
} from 'native-base';

import { PTSansText } from '../components/StyledText';

const Profile = ({ profile }) => {

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
        {
          profile.pictureUrl !== '' ?
            <Thumbnail source={{ uri: profile.pictureUrl }} />
          :
            <Thumbnail source={require('../../assets/images/icon.png')} />
        }
      </Left>
      <Body>
        <PTSansText style={styles.name}>{profile.formattedName}</PTSansText>
        <PTSansText>{profile.headline}</PTSansText>
        <PTSansText>
          <Icon name='pin' style={styles.location} /> {profile.location.name}
        </PTSansText>
      </Body>
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

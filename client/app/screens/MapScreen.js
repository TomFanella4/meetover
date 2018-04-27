import React from 'react';
import { Platform, StyleSheet, Image } from 'react-native';
import { Permissions, MapView, Marker, Callout } from 'expo';
import { connect } from 'react-redux';
import { Ionicons } from '@expo/vector-icons';
import { Spinner, View } from 'native-base';

import { PTSansText } from '../components/StyledText';
import { ProfileImage } from '../components/ProfileImage';
import IsSearchingBar from '../components/IsSearchingBar';
import Colors from '../constants/Colors';

class MapScreen extends React.Component {
  static navigationOptions = {
    title: 'Map',
  };

  state = {
    location: null
  }

  _viewProfile(match) {
    const { profile, matchStatus } = match;
    // TODO clean up data flow
    this.props.navigation.navigate('RequestScreen', { profile: {...profile, ...matchStatus} });
  }

  render() {
    const { matches } = this.props;

    const matchMarkers = matches.map(match => (
      match.location &&
      <MapView.Marker coordinate={match.location} key={match.profile.id}>
        <MapView.Callout
          style={styles.mapViewCallout}
          onPress={() => this._viewProfile(match)}
        >
          <ProfileImage pictureUrl={match.profile.pictureUrl} />
          <View style={styles.mapViewCalloutSection}>
            <PTSansText style={styles.name}>
              {match.profile.formattedName}
            </PTSansText>
            <PTSansText style={styles.headline}>
              {match.profile.headline.substring(0,25)}
              {match.profile.headline.length >= 25 && '...'}
            </PTSansText>
          </View>
          <View style={styles.mapViewCalloutSection}>
            <Ionicons
              name={Platform.OS === 'ios' ? 'ios-information-circle' : 'md-information-circle'}
              size={30}
            />
          </View>
        </MapView.Callout>
      </MapView.Marker>
    ));

    return (
      <View style={styles.container}>
        <IsSearchingBar />
        {
          this.state.location ?
            <MapView
              style={styles.mapView}
              initialRegion={this.state.location}
              showsUserLocation={true}
            >
              {matchMarkers}
            </MapView>
          :
            <Spinner color={Colors.tintColor} />
        }
      </View>
    );
  }

  componentDidMount() {
    this.updateLocation();
  }

  async updateLocation() {
    const { status } = await Permissions.askAsync(Permissions.LOCATION);
    if (status === 'granted') {
      const location = await Expo.Location.getCurrentPositionAsync({})
      .catch(err => console.error(err));

      this.setState({
        location: {
          latitude: location.coords.latitude,
          longitude: location.coords.longitude,
          latitudeDelta: 0.0922,
          longitudeDelta: 0.0421,
        }
      });
    } else {
      throw new Error('Location permission not granted');
    }
  }
};

const mapStateToProps = state => ({
  matches: state.matchList.matches
});

export default connect(
  mapStateToProps,
)(MapScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  mapView: {
    flex: 1,
  },
  mapViewCallout: {
    flexDirection: 'row',
  },
  mapViewCalloutSection: {
    paddingLeft: 10,
    justifyContent: 'center',
  },
  name: {
    fontSize: 20,
  },
  headline: {
    fontSize: 12,
  },
});

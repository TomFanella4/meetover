import React from 'react';
import { ScrollView, StyleSheet } from 'react-native';
// import { connect } from 'react-redux';

import Colors from '../constants/Colors';
import { PTSansText } from '../components/StyledText'

export default class ListScreen extends React.Component {
  static navigationOptions = {
    title: 'List',
    headerStyle: {
      backgroundColor: Colors.tintColor,
    },
    headerTintColor: Colors.header,
  };

  render() {
    return (
      <ScrollView style={styles.container}>
        <PTSansText>List Screen</PTSansText>
      </ScrollView>
    );
  }
};

// TODO: Implement ListScreen functionality
// const mapStateToProps = state => ({
//
// });
//
// const mapDispatchToProps = {
//
// };
//
// export default connect(
//   mapStateToProps,
//   mapDispatchToProps
// )(ListScreen);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingTop: 15,
    backgroundColor: '#fff',
  }
});

import React from 'react';
import { ScrollView, StyleSheet } from 'react-native';

// import { bindActionCreators } from 'redux'
// import { connect } from 'react-redux';
//
// import * as Actions from '../actions';

import { PTSansText } from '../components/StyledText'

export default class ListScreen extends React.Component {
  static navigationOptions = {
    title: 'List',
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
// const mapStateToProps = (state) => ({
//
// });
//
// const mapDispatchToProps = (dispatch) => (
//   bindActionCreators(Actions, dispatch)
// );
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

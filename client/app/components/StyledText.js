import React from 'react';
import { Text } from 'native-base';

export class PTSansText extends React.Component {
  render() {
    return <Text {...this.props} style={[this.props.style, { fontFamily: 'pt-sans' }]} />;
  }
}

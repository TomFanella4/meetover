import React from 'react';
import { Thumbnail } from 'native-base';

export const ProfileImage = ({ pictureUrl, style }) => (
  pictureUrl !== '' ?
    <Thumbnail style={style} source={{ uri: pictureUrl }} />
  :
    <Thumbnail style={style} source={require('../../assets/images/icon.png')} />
);

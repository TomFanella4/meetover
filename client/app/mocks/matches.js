import { times, random } from 'lodash';

const profile = {
  "currentShare": {
    "attribution": {"share": {
      "author": {
        "firstName": "MaryAnn",
        "id": "bxQrUhmst0",
        "lastName": "Gibney"
      },
      "comment": "Big thanks to the always vibrant Geoff Nyheim and Amazon Web Services colleagues Roshni Joshi, Michael Dowling ☁, Peter Tannenwald, Robert Heitzler, and @shane for coaching and mentoring DePaul University Center for Sales Leadership students on career and resume skills this afternoon! Plus DePaul alum Jeremy Paul for sharing your experience! #CSLyourself",
      "id": "s6362854320802713600"
    }},
    "author": {
      "firstName": "Krutarth",
      "id": "HJSNGIIRCj",
      "lastName": "Rao"
    },
    "comment": "Purdue University ",
    "id": "s6361684091376668672",
    "source": {"serviceProvider": {"name": "FLAGSHIP"}},
    "timestamp": 1516743681760,
    "visibility": {"code": "anyone"}
  },
  "emailAddress": "krk91@outlook.com",
  "firstName": "Krutarth",
  "formattedName": "Krutarth Rao",
  "headline": "Software Engineering Intern at Aruba, a Hewlett Packard Enterprise company",
  "id": "HJSNGIIRCj",
  "industry": "Computer Software",
  "lastName": "Rao",
  "location": {
    "country": {"code": "us"},
    "name": "United States"
  },
  "numConnections": 411,
  "numConnectionsCapped": false,
  "pictureUrl": "https://media.licdn.com/mpr/mprx/0_1lpYsrJfPnDiMc-pBPEPZ-hfnNriM6Cg9vEYxXza1zwCM91rsle0tT4fvFjCMnPlnbXKlGc70KuGVGW7JGJCtXs_sKu_VGBjsGJgML7mllMfrquT9P71JlTxVLnPPGh21TgtORGj196",
  "positions": {
    "_total": 1,
    "values": [{
      "company": {
        "id": 3846,
        "industry": "Higher Education",
        "name": "Purdue University",
        "size": "10001+",
        "type": "Educational"
      },
      "id": 827836295,
      "isCurrent": true,
      "location": {"name": "West Lafayette, Indiana"},
      "startDate": {
        "month": 5,
        "year": 2016
      },
      "summary": "Current:\n\nDesigning a supply chain management system using a blockchain oriented approach. Implemented a prototype using the Hyperledger project by IBM and discovering interesting use cases and research aspects for the system.\nVisit https://freedom.cs.purdue.edu/ for more info.\n\nPrevious Project:\n\nWorked in a team to design a communication protocol using cryptography primitives and cryptocurrencies such as Bitcoin and Ether. Some key knowledge I acquired for the project includes cryptographic hashing, watermarking (robust, fragile and semi-fragile), oblivious transfer protocol, blockchain primitives, bitcoin scripts, smart contracts, traitor tracing and other related fields in order to design the intended system and provide adequate applications to it.",
      "title": "Research Assistant"
    }]
  },
  "summary": "Senior (Graduation May 2018) at Purdue university pursuing a BSc in Computer Science looking to apply my knowledge in the field of software development, system security and blockchain solutions through full-time opportunities."
};

export default matchesMock = times(10, () => ({
  ...profile,
  "location": {
    "latitude": 40.424649 + random(-0.01, 0.01, true),
    "longitude": -86.911571 + random(-0.01, 0.01, true),
  },
}));

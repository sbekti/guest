import * as React from 'react';
import {
  Container,
  Stack,
  Text,
  Heading,
  useColorModeValue,
} from '@chakra-ui/react';

const sectionText = text => {
  return (
    <Text
      color={'green.500'}
      textTransform={'uppercase'}
      fontWeight={800}
      fontSize={'sm'}
      letterSpacing={1.1}
    >
      {text}
    </Text>
  );
};

export const TermsPage = () => {
  return (
    <Container maxW={'7xl'} p="12">
      <Stack>
        <Heading
          color={useColorModeValue('gray.700', 'white')}
          fontSize={'2xl'}
          fontFamily={'body'}
        >
          Guest Wi-Fi Wireless Networking Acceptable Use Policy
        </Heading>
        <Text color={'gray.500'}>
          We are offering this guest Wi-Fi wireless Internet service (the
          “Service”) according to this Guest Wi-Fi Wireless Networking
          Acceptable Use Policy (the “Policy”) as a free, non-public service to
          its visitors for the duration of their official visits. All users of
          this Service must agree to the terms of this Policy by checking the
          box on the registration page. We do not guarantee the Service or
          specific rates of speed. We also have no control over information
          obtained through the Internet and cannot be held responsible for its
          content or accuracy. Use of the service is subject to the user’s own
          risk.We reserve the right to remove, block, filter, or restrict by any
          other means any material that, in our sole discretion, may be illegal,
          may subject us to liability, or may violate this Policy. We may
          cooperate with legal authorities and/or third parties in the
          investigation of any suspected or alleged crime or civil wrong.
          Violations of this Policy may result in the suspension or termination
          of access to the Service or other resources, or other actions as
          detailed below.
        </Text>
        {sectionText('Responsibilities of Service Users')}
        <Text color={'gray.500'}>
          Users are responsible for ensuring they are running up-to-date
          anti-virus software on their wireless devices. Users must be aware
          that, as they connect their devices to the Internet through the
          Service, they expose their devices to: worms, viruses, Trojan horses,
          denial-of-service attacks, intrusions, packet-sniffing, and other
          abuses by third-parties. Users must respect all copyrights.
          Downloading or sharing copyrighted materials is strictly prohibited.
          The running of programs, services, systems, processes, or servers by a
          single user or group of users that may substantially degrade network
          performance or accessibility will not be allowed. Electronic chain
          letters and mail bombs are prohibited. Connecting to "Peer to Peer"
          file sharing networks or downloading large files, such as CD ISO
          images, is also prohibited.Accessing another person's computer,
          computer account, files, or data without permission is prohibited.
          Attempting to circumvent or subvert system or network security
          measures is prohibited. Creating or running programs that are designed
          to identify security loopholes, to decrypt intentionally secured data,
          or to gain unauthorized access to any system is prohibited. Using any
          means to decode or otherwise obtain restricted passwords or access
          control information is prohibited. Forging the identity of a user or
          machine in an electronic communication is prohibited. Saturating
          network or computer resources to the exclusion of another's use, for
          example, by overloading the network with traffic such as emails or
          legitimate (file backup or archive) or malicious (denial of service
          attack) activity, is prohibited. Users understand that wireless
          Internet access is inherently not secure, and users should adopt
          appropriate security measures when using the Service. We highly
          discourage users from conducting confidential transactions (such as
          online banking, credit card transactions, etc.) over any wireless
          network, including this Service. Users are responsible for the
          security of their own devices.
        </Text>
        {sectionText('Limitations of Wireless Network Access')}
        <Text color={'gray.500'}>
          We are not liable for any damage, undesired resource usage, or
          detrimental effects that may occur to a user's device and/or software
          while the user’s device is attached to the Service. The user is
          responsible for any actions taken from his or her device, whether
          intentional or unintentional, that damage or otherwise affect other
          devices or users of the Service. The user hereby releases the Company
          from liability for any loss, damage, security infringement, or injury
          which the user may sustain as a result of being allowed access to the
          Service. The user agrees to be solely responsible for any such loss,
          infringement, damage, or injury.
        </Text>
        {sectionText('Terms of Service')}
        <Text color={'gray.500'}>
          By checking the box in the registration page, the user agrees to
          comply with and to be legally bound by the terms of this Policy. If
          this Policy or any terms of the Service are unacceptable or become
          unacceptable to the user, the user's only right shall be to terminate
          his or her use of the Service.
        </Text>
        {sectionText('Lawful Use')}
        <Text color={'gray.500'}>
          The Service may only be used for lawful purposes and in a manner which
          we believe to be consistent with the rights of other users. The
          Service shall not be used in a manner which would violate any law or
          infringe any copyright, trademark, trade secret, right of publicity,
          privacy right, or any other right of any person or entity. The Service
          shall not be used for the purpose of accessing, transmitting, or
          storing material which is considered obscene, libelous or defamatory.
          Illegal acts may subject users to prosecution by local, state,
          federal, or international authorities. We may bring legal action to
          enjoin violations of this Policy and/or to collect damages, if any,
          caused by violations.
        </Text>
        {sectionText(
          'The user specifically agrees to the following conditions'
        )}
        <Text color={'gray.500'}>
          The user will use the Service only as permitted by applicable local,
          state, federal, and International laws. The user will refrain from any
          actions that we consider to be negligent or malicious. The user will
          not send email containing viruses or other malicious or damaging
          software. The user will run appropriate anti-virus software to remove
          such damaging software from his or her computer. The user will not
          access web sites which contain material that is grossly offensive to
          us, including clear expressions of bigotry, racism, or hatred. The
          user will not access web sites which contain material that defames,
          abuses, or threatens others.
        </Text>
        {sectionText('Changes to Service')}
        <Text color={'gray.500'}>
          We reserve the right to change the Service offered, the features of
          the Service offered, the terms of this Policy, or its system without
          notice to the user.
        </Text>
      </Stack>
      <Stack mt={6}>
        <Text fontWeight={600}>
          By registering to the Service, you acknowledge that you understand and
          agree to this Policy.
        </Text>
      </Stack>
    </Container>
  );
};

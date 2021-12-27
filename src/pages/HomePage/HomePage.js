import { useEffect, useState } from 'react';
import {
  Flex,
  Stack,
  Heading,
  Text,
  useColorModeValue,
} from '@chakra-ui/react';
import { RegisterCard } from '../../components/RegisterCard';
import { SuccessCard } from '../../components/SuccessCard';

export const HomePage = () => {
  const [showSuccess, setShowSuccess] = useState(false);

  useEffect(() => {
    const appHeight = () => {
      const doc = document.documentElement;
      doc.style.setProperty('--app-height', `${window.innerHeight}px`);
    }
    window.addEventListener('resize', appHeight);
    appHeight();
  }, []);

  const handleRegisterSuccess = (result) => {
    setShowSuccess(true);
  }

  const componentToDisplay = () => {
    if (showSuccess) {
      return (
        <Stack spacing={8} mx={'auto'} maxW={'lg'} py={12} px={6}>
          <Stack align={'center'}>
            <SuccessCard onBackToHomeClick={() => setShowSuccess(false)} />
          </Stack>
        </Stack>
      );
    } else {
      return (
        <Stack spacing={8} mx={'auto'} maxW={'lg'} py={12} px={6}>
          <Stack align={'center'}>
            <Heading fontSize={'4xl'}>Bektinet</Heading>
            <Text fontSize={'lg'} color={'gray.600'}>
              Register to access the guest Wi-Fi
            </Text>
          </Stack>
          <RegisterCard onRegisterSuccess={handleRegisterSuccess} />
        </Stack>
      );
    }
  }

  return (
    <Flex
      minH={'var(--app-height)'}
      align={'center'}
      justify={'center'}
      bg={useColorModeValue('gray.50', 'gray.800')}>
      { componentToDisplay() }
    </Flex>
  );
};

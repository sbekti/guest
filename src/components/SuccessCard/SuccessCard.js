import {
  Stack,
  Heading,
  Text,
  Button,
  Icon,
  useColorModeValue,
  createIcon,
} from '@chakra-ui/react';
  
export const SuccessCard = ({ onBackToHomeClick }) => {
  return (
    <Stack
      bg={useColorModeValue('white', 'gray.700')}
      borderWidth='1px'
      rounded={'lg'}
      p={4}
      spacing={8}
      align={'center'}>
      <Icon as={NotificationIcon} w={24} h={24} />
      <Stack align={'center'} spacing={2}>
        <Heading
          fontSize={'3xl'}
          color={useColorModeValue('gray.800', 'gray.200')}>
          Success!
        </Heading>
        <Text fontSize={'lg'} color={'gray.500'} align={'center'}>
          Please check your inbox. We have sent an email containing the login information.
        </Text>
      </Stack>
      <Stack spacing={4} direction={{ base: 'column', md: 'row' }} w={'full'}>
        <Button
          bg={'blue.400'}
          color={'white'}
          flex={'1 0 auto'}
          _hover={{ bg: 'blue.500' }}
          _focus={{ bg: 'blue.500' }}
          onClick={onBackToHomeClick}>
          Back to Home
        </Button>
      </Stack>
    </Stack>
  );
}

const NotificationIcon = createIcon({
  displayName: 'Notification',
  viewBox: '0 0 128 128',
  path: (
    <g id='Notification'>
      <rect
        className='cls-1'
        x='1'
        y='45'
        fill={'#fbcc88'}
        width='108'
        height='82'
      />
      <circle className='cls-2' fill={'#8cdd79'} cx='105' cy='86' r='22' />
      <rect
        className='cls-3'
        fill={'#f6b756'}
        x='1'
        y='122'
        width='108'
        height='5'
      />
      <path
        className='cls-4'
        fill={'#7ece67'}
        d='M105,108A22,22,0,0,1,83.09,84a22,22,0,0,0,43.82,0A22,22,0,0,1,105,108Z'
      />
      <path
        fill={'#f6b756'}
        className='cls-3'
        d='M109,107.63v4A22,22,0,0,1,83.09,88,22,22,0,0,0,109,107.63Z'
      />
      <path
        className='cls-5'
        fill={'#d6ac90'}
        d='M93,30l16,15L65.91,84.9a16,16,0,0,1-21.82,0L1,45,17,30Z'
      />
      <path
        className='cls-6'
        fill={'#cba07a'}
        d='M109,45,65.91,84.9a16,16,0,0,1-21.82,0L1,45l2.68-2.52c43.4,40.19,41.54,39.08,45.46,40.6A16,16,0,0,0,65.91,79.9l40.41-37.42Z'
      />
      <path
        className='cls-7'
        fill={'#dde1e8'}
        d='M93,1V59.82L65.91,84.9a16,16,0,0,1-16.77,3.18C45.42,86.64,47,87.6,17,59.82V1Z'
      />
      <path
        className='cls-3'
        fill={'#f6b756'}
        d='M46.09,86.73,3,127H1v-1c6-5.62-1.26,1.17,43.7-40.78A1,1,0,0,1,46.09,86.73Z'
      />
      <path
        className='cls-3'
        fill={'#f6b756'}
        d='M109,126v1h-2L63.91,86.73a1,1,0,0,1,1.39-1.49C111,127.85,103.11,120.51,109,126Z'
      />
      <path
        className='cls-8'
        fill={'#c7cdd8'}
        d='M93,54.81v5L65.91,84.9a16,16,0,0,1-16.77,3.18C45.42,86.64,47,87.6,17,59.82v-5L44.09,79.9a16,16,0,0,0,21.82,0Z'
      />
      <path
        className='cls-9'
        fill={'#fff'}
        d='M101,95c-.59,0-.08.34-8.72-8.3a1,1,0,0,1,1.44-1.44L101,92.56l15.28-15.28a1,1,0,0,1,1.44,1.44C100.21,96.23,101.6,95,101,95Z'
      />
    </g>
  ),
});

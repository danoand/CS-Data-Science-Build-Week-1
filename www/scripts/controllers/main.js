/**
 *
 * appCtrl
 *
 */

angular
    .module('homer')
    .controller('appCtrl', appCtrl)
    .controller('dashCtrl', dashCtrl);

function appCtrl($http, $scope) {
}

function dashCtrl($http, $scope) {
    $scope.info = {};

    $http({
        method: 'GET',
        url: '/getModelInfo'
      }).then(function successCallback(response) {
          console.log('success response: ' + JSON.stringify(response));
          $scope.info = response.data.content;
        }, function errorCallback(response) {
            console.log('error response: ' + JSON.stringify(response));
        });
}

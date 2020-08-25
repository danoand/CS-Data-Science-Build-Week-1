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

function dashCtrl($http, $scope, growl) {
    $scope.info             = {};
    $scope.pred             = {};
    $scope.pred.havepred    = false;
    $scope.pypred           = {};
    $scope.pypred.havepred  = false;

    $http({
        method: 'GET',
        url: '/getModelInfo'
      }).then(function successCallback(response) {
          console.log('success response: ' + JSON.stringify(response));
          $scope.info = response.data.content;
        }, function errorCallback(response) {
            console.log('error response: ' + JSON.stringify(response));
        });

    $scope.getPred = function() {
        $scope.pred.havepred = false;

        $http({
            method: 'POST',
            data: {"text": $scope.spam_text},
            url: '/getModelPred'
          }).then(function successCallback(response) {
              console.log('success response: ' + JSON.stringify(response));
              $scope.pred = response.data.content;
            }, function errorCallback(response) {
                console.log('error response: ' + JSON.stringify(response));
                growl.warning(response.data.msg, {ttl: 2500});
            });

        $http({
            method: 'POST',
            data: {"text": $scope.spam_text},
            url: '/getPyModelPred'
            }).then(function successCallback(rsp) {
                console.log('success response: ' + JSON.stringify(rsp));
                $scope.pypred = rsp.data.content;
            }, function errorCallback(rsp) {
                console.log('error response: ' + JSON.stringify(rsp));
                growl.warning(rsp.data.msg, {ttl: 2500});
            });
    };
}
